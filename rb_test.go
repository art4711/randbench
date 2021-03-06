package randbench

import (
	"io"
	"math/rand"
	"testing"

	"github.com/art4711/unpredictable"
)

func do(b *testing.B, src rand.Source) {
	b.ReportAllocs()
	r := rand.New(src)
	for i := 0; i < b.N; i++ {
		_ = r.Int63()
	}
}

func BenchmarkMathRand(b *testing.B) {
	do(b, rand.NewSource(1))
}

func BenchmarkCryptoRand(b *testing.B) {
	do(b, NewCryptoRand())
}

func BenchmarkLibcRand(b *testing.B) {
	do(b, NewLibcRand())
}

func BenchmarkLibcBufRand(b *testing.B) {
	do(b, NewLibcBufRand())
}

func BenchmarkCryptoBufRand(b *testing.B) {
	do(b, NewCryptoBufRand())
}

func BenchmarkCOverhead(b *testing.B) {
	do(b, NewCOverhead())
}

func BenchmarkUnpredictable(b *testing.B) {
	do(b, unpredictable.NewMathRandSource())
}

func BenchmarkCryptoCastRand(b *testing.B) {
	do(b, NewCryptoCastRand())
}

const sz = 2 * 1024 * 1024

var scratch [sz]byte

func doread(b *testing.B, src io.Reader) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := src.Read(scratch[:])
		if err != nil {
			b.Fatal(err)
		}
		b.SetBytes(sz)
	}
}

func BenchmarkReadMathRand(b *testing.B) {
	doread(b, rand.New(rand.NewSource(1)))
}

func BenchmarkReadUnpredictable(b *testing.B) {
	doread(b, unpredictable.NewReader())
}
