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

### COverhead ###

I'm curious, how big is the overhead of calling into C anyway?

### Unpredictable ###

My own attempt at implementing parts of OpenBSD arc4random in Go.
See: https://github.com/art4711/unpredictable

### CryptoCast ###

How much overhead is encoding/binary.Read anyway?

## The results ##

    BenchmarkMathRand-4     	200000000	         9.86 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCryptoRand-4   	 2000000	       790 ns/op	      16 B/op	       2 allocs/op
    BenchmarkLibcRand-4     	10000000	       227 ns/op	       0 B/op	       0 allocs/op
    BenchmarkLibcBufRand-4  	50000000	        36.4 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCryptoBufRand-4	 3000000	       573 ns/op	       8 B/op	       0 allocs/op
    BenchmarkCOverhead-4    	10000000	       196 ns/op	       0 B/op	       0 allocs/op
    BenchmarkUnpredictable-4	30000000	        42.3 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCryptoCastRand-4	 2000000	       724 ns/op	       8 B/op	       1 allocs/op

We can see that the default non-random number generator takes 10ns per
op. crypto/rand is 80 times slower. That's bad, although I suspect
that it has more to do with encoding/binary than the generator itself.
For some weird reason Go decided to not have the `crypto/rand`
generator pluggable into the nice framework from `math/rand` so we end
up converting between integers and byte arrays back and forth.

The interesting part is LibcRand and LibcBufRand. LibcRand is 23 times
slower which isn't very bad, but it's a magnitude worse than I'd like
it to be. But I had a strong suspicion that this is more due to three
things:

 - Go calling into C has non-trivial overhead. I'm pretty sure we're
   messing with the stack and poke into the scheduler.
 - The old version of `arc4random` on my test system (MacOS) has
   very expensive fork detection code (it calls `getpid`).
 - `arc4random` locks for every call. math/rand doesn't.

So to test that theory let's buffer the numbers that libc returns.
LibcBufRand does that. And suddenly we're just 3.6x slower. This is
perfectly acceptable for vastly improved quality.

Just to make sure that the buffer doesn't do something magical,
CryptoBufRand implements the same buffer around the same source as
CryptoRand, but as we can see, it doesn't magically become much
better.

COverhead is there to actually see the overhead of calling into
C. Holy crap.

Unpredictable is my implementation of arc4random without the
overhead of calling into C. 4x slower, not bad at all.

CryptoCast is there to see how much encoding/binary.Read adds to
reading from crypto/rand. Quite a bit. Loos like 70ns per op. Still
doesn't put crypto/rand anywhere near the other sources though.

## The results on Go 1.9 ##

    BenchmarkMathRand-4         	200000000	         7.97 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCryptoRand-4       	 2000000	       729 ns/op	      16 B/op	       2 allocs/op
    BenchmarkLibcRand-4         	 5000000	       359 ns/op	       0 B/op	       0 allocs/op
    BenchmarkLibcBufRand-4      	50000000	        33.9 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCryptoBufRand-4    	 3000000	       599 ns/op	       8 B/op	       0 allocs/op
    BenchmarkCOverhead-4        	20000000	        77.5 ns/op	       0 B/op	       0 allocs/op
    BenchmarkUnpredictable-4    	50000000	        26.2 ns/op	       0 B/op	       0 allocs/op
    BenchmarkCryptoCastRand-4   	 2000000	       706 ns/op	       8 B/op	       1 allocs/op

Unpredictable is slowly catching up to math/rand.
