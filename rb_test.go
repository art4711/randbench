package rb

import (
	"testing"
	"math/rand"
)

func BenchmarkMathRand(b *testing.B) {
	src := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		_ = src.Int63()
	}
}

func BenchmarkCryptoRand(b *testing.B) {
	src := rand.New(NewCryptoRand())
	for i := 0; i < b.N; i++ {
		_ = src.Int63()
	}
}

func BenchmarkLibcRand(b *testing.B) {
	src := rand.New(NewLibcRand())
	for i := 0; i < b.N; i++ {
		_ = src.Int63()
	}
}

func BenchmarkLibcBufRand(b *testing.B) {
	src := rand.New(NewLibcBufRand())
	for i := 0; i < b.N; i++ {
		_ = src.Int63()
	}
}
