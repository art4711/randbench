# Why do standard libraries of modern languages always implement bad random? #

For some reason modern computer science has decided that "random
numbers" should mean "perfectly predictable numbers" with the added
requirement "each run of a program should produce the exact same
sequence unless you do some unclear operation that doesn't improve
much" and usually "forked children should produce the exact same
sequences as each other".

The normal argument for this is that people don't actually want
properly random numbers (multiple security failures caused by
predictable random number generators show otherwise), most people want
repeatable sequences (research by OpenBSD showed that 90-something% of
all programs in the ports tree go out of their way to avoid them). And
finally, because decent random number generators are slow.

Let's test that last statement. How bad can it be?

## The benchmarks ##

All tests follow the same structure:

    r := rand.New(SOURCE)
    for i := 0; i < b.N; i++ {
    	_ = r.Int63()
    }

Let's keep it simple. `SOURCE` is the part that varies between the
tests.

## Our sources ##

All of our sources implement the `math/rand.Source` interface. But
since our sources are not seedable, they will panic if `Seed` is
called. The only other function they need to implement is `Int63`. Why
63 bits and not 64? Beats me, but apparently this is very important
because that restriction is heavily documented in function names and
everywhere else. Negative numbers are bad and so are unsigned 64 bit
numbers. I don't know why, it actually is an extra restriction on our
sources because they need to mask the high bits out while the default
Go source doesn't. This smells like bug that was documented away.

### MathRand ###

This tests the default (not at all)"random" number generator in the Go
standard library. Source is `rand.NewSource(1)`, I seed to 1 because
that's the default seed.

### CryptoRand ###

`crypto/rand.Reader` is read through `encoding/binary.Read` into an
int64 and high bit is masked out. I suspect that a lot of the overhead
is in `encoding/binary`, but I've made no effort to fix it.

### LibcRand ###

We assume we're on a sane libc that implements `arc4random_buf`. We
read into a `uint64_t` and mask out the high bit.

### LibcBufRand ###

Same as above, but we request 256 values from libc in one go and then
return from the buffer.

### CryptoBufRand ###

Like CryptoRand, but with a buffer like in LibcBufRand.

## The results ##

    BenchmarkMathRand-4     	100000000	        10.0 ns/op
    BenchmarkCryptoRand-4   	 2000000	       840 ns/op
    BenchmarkLibcRand-4     	10000000	       240 ns/op
    BenchmarkLibcBufRand-4  	50000000	        36.9 ns/op
    BenchmarkCryptoBufRand-4	 2000000	       613 ns/op

We can see that the default non-random number generator takes 10ns per
op. crypto/rand is 84 times slower. That's bad, although I suspect
that it has more to do with encoding/binary than the generator itself.
For some weird reason Go decided to not have the `crypto/rand`
generator pluggable into the nice framework from `math/rand` so we end
up converting between integers and byte arrays back and forth.

The interesting part is LibcRand and LibcBufRand. LibcRand is 24 times
slower which isn't very bad, but it's a magnitude worse than I'd like
it to be. But I had a strong suspicion that this is more due to three
things:

 - Go has a different function calling convention than C and
   translating the arguments between the two is expensive.
 - The old version of `arc4random` on my test system (MacOS) has
   very expensive fork detection code (it calls `getpid`).
 - `arc4random` locks for every call. math/rand doesn't.

So to test that theory let's buffer the numbers that libc returns.
LibcBufRand does that. And suddenly we're just 3.7x slower. This is
perfectly acceptable for vastly improved quality.

Just to make sure that the buffer doesn't do something magical,
CryptoBufRand implements the same buffer around the same source as
CryptoRand, but as we can see, it doesn't magically become much
better.

## Future tests ##

### modern arc4random ###

More modern `arc4random` code has much better fork detection, so at
least that part can be improved. Also, the algorithm has changed, I'd
really want to see how that performs. Maybe even implement the same
algorithm in Go, just for comparison.

### better reading from crypto/rand ###

I'd really want to get rid of the overhead of `encoding/binary`. It's
fun to poke fun at how a server language is bad at encoding and
decoding data, but it would make sense to remove that stupid overhead
to test the actual generator.
