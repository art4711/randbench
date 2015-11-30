package randbench

// #include <stdint.h>
// #include <stdlib.h>
// static int64_t r64(void) { uint64_t r; arc4random_buf(&r, sizeof(r)); return r & 0x7fffffffffffffffLL; }
import "C"

import (
	badrand "math/rand"
)

type libcsrc struct {
	genericSrc
}

func (s *libcsrc) Int63() int64 {
	return int64(C.r64())
}

func NewLibcRand() badrand.Source {
	return &libcsrc{}
}
