package rb

import (
	goodrand "crypto/rand"
	"encoding/binary"
	badrand "math/rand"
)

type src struct {
	genericSrc
}

func (s *src) Int63() int64 {
	var i int64
	err := binary.Read(goodrand.Reader, binary.LittleEndian, &i)
	if err != nil {
		panic(err)
	}
	return i << 1 >> 1
}

func NewCryptoRand() badrand.Source {
	return &src{}
}
