package randbench

import (
	goodrand "crypto/rand"
	badrand "math/rand"
	"unsafe"
)

type cryptocastsrc struct {
	genericSrc
}

func (s *cryptocastsrc) Int63() int64 {
	var b [8]byte	
	_, err := goodrand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return *(*int64)(unsafe.Pointer(&b[0])) & 0x7fffffffffffffff
}

func NewCryptoCastRand() badrand.Source {
	return &cryptocastsrc{}
}
