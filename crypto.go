package randbench

import (
	goodrand "crypto/rand"
	"encoding/binary"
	badrand "math/rand"
)

type cryptosrc struct {
	genericSrc
}

func (s *cryptosrc) Int63() int64 {
	var i int64
	err := binary.Read(goodrand.Reader, binary.LittleEndian, &i)
	if err != nil {
		panic(err)
	}
	return i & 0x7fffffffffffffff
}

func NewCryptoRand() badrand.Source {
	return &cryptosrc{}
}
