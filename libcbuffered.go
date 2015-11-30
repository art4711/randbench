package randbench

// #include <stdint.h>
// #include <stdlib.h>
import "C"

import (
	badrand "math/rand"
	"unsafe"
)

const bufsz = 256

type libcbufsrc struct {
	genericSrc
	buf [bufsz]int64
	bp  int
}

func (s *libcbufsrc) stir() {
	C.arc4random_buf(unsafe.Pointer(&s.buf), C.size_t(unsafe.Sizeof(s.buf)))
}

func (s *libcbufsrc) Int63() int64 {
	if s.bp == 0 {
		s.stir()
	}
	r := s.buf[s.bp]
	s.bp = (s.bp + 1) % bufsz
	return r & 0x7fffffffffffffff
}

func NewLibcBufRand() badrand.Source {
	return &libcbufsrc{}
}
