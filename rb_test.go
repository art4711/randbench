package randbench

import (
	"math/rand"
	"testing"
)

func do(b *testing.B, src rand.Source) {
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
