package rb

import (
	badrand "math/rand"
	goodrand "crypto/rand"
	"encoding/binary"
)

type src struct {
}

func (s *src)Int63() int64 {
	var i int64
	err := binary.Read(goodrand.Reader, binary.LittleEndian, &i)
	if err != nil {
		panic(err)
	}
	return i << 1 >> 1
}

func (s *src)Seed(seed int64) {
	panic("if you seed this source you need to reconsider the choices you made in your life that led you to this")
}

func NewCryptoRand() badrand.Source {
	return &src{}
}
