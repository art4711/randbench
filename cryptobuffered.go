package rb

import (
	goodrand "crypto/rand"
	"encoding/binary"
	badrand "math/rand"
)

const cbufsz = 256

type cryptobufferedsrc struct {
	genericSrc
	buf [cbufsz]int64
	bp int
}

func (s *cryptobufferedsrc) stir() {
	err := binary.Read(goodrand.Reader, binary.LittleEndian, &s.buf)
	if err != nil {
		panic(err)
	}
}

func (s *cryptobufferedsrc) Int63() int64 {
	if s.bp == 0 {
		s.stir()
	}
	r := s.buf[s.bp]
	s.bp = (s.bp + 1) % cbufsz
	return r << 1 >> 1
}

func NewCryptoBufRand() badrand.Source {
	return &cryptobufferedsrc{}
}
