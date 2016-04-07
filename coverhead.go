package randbench

// #include <stdint.h>
// #include <stdlib.h>
// static int64_t r17(void) { return 17; }
import "C"

import (
	badrand "math/rand"
)

type co struct {
}

func (s co) Int63() int64 {
	return int64(C.r17())
}

func (s co) Seed(a int64) {
	panic("no")
}

func NewCOverhead() badrand.Source {
	return co{}
}
