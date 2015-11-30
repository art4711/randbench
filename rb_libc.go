package rb

// #include <stdint.h>
// #include <stdlib.h>
// static int64_t r64(void) { uint64_t r; arc4random_buf(&r, sizeof(r)); return r & 0x7fffffffffffffffLL; }
import "C"

import (
	badrand "math/rand"	
)

type libcsrc int

func (s libcsrc)Int63() int64 {
	return int64(C.r64())
}

func (s libcsrc)Seed(seed int64) {
	panic("if you seed this source you need to reconsider the choices you made in your life that led you to this")
}

func NewLibcRand() badrand.Source {
	return libcsrc(0)
}
