package vm

import (
	"bytes"
	"encoding/binary"
	"math"
)

var enc = binary.LittleEndian

type Generator struct {
	buffer bytes.Buffer
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (s *Generator) Op(op Op) *Generator {
	s.buffer.WriteByte(byte(op))
	return s
}

func (s *Generator) Uint16(v uint16) *Generator {
	var data [2]byte
	enc.PutUint16(data[:], v)
	s.buffer.Write(data[:])
	return s
}

func (s *Generator) String(v string) *Generator {
	if len(v) > math.MaxUint16 {
		panic("string is too long")
	}

	s.Uint16(uint16(len(v)))
	s.buffer.Write([]byte(v))
	return s
}

func (s *Generator) Int64(v int64) *Generator {
	var data [8]byte
	enc.PutUint64(data[:], uint64(v))
	s.buffer.Write(data[:])
	return s
}
