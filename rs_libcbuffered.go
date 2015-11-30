package rb

// #include <stdint.h>
// #include <stdlib.h>
import "C"

import (
	badrand "math/rand"
	"unsafe"
)

const bufsz = 256

type libcbufsrc struct {
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
	return r << 1 >> 1
}

func (s *libcbufsrc) Seed(seed int64) {
	panic("if you seed this source you need to reconsider the choices you made in your life that led you to this")
}

func NewLibcBufRand() badrand.Source {
	return &libcbufsrc{}
}
