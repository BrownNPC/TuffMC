package ds

import (
	"encoding/binary"
	"math"
)

// A single-precision 32-bit IEEE 754 floating point number
func EncodeFloat(v float32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], math.Float32bits(v))
	return b[:]
}

// A single-precision 32-bit IEEE 754 floating point number
func DecodeFloat(b []byte) float32 {
	v := binary.BigEndian.Uint32(b)
	return math.Float32frombits(v)
}

// A double-precision 64-bit IEEE 754 floating point number
func EncodeDouble(v float64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], math.Float64bits(v))
	return b[:]
}

// A double-precision 64-bit IEEE 754 floating point number
func DecodeDouble(b []byte) float64 {
	v := binary.BigEndian.Uint64(b)
	return math.Float64frombits(v)
}
