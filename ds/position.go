package ds

import "encoding/binary"

// https://minecraft.wiki/w/Data_types?oldid=2767030#Position
// 64-bit value split in to three parts
// x: 26 MSBs
// z: 26 LSBs
// y: 12 bits between them
// Encoded as followed:
//
//	((x & 0x3FFFFFF) << 38) | ((y & 0xFFF) << 26) | (z & 0x3FFFFFF)
func EncodePosition(X, Y, Z int32) []byte {
	x := uint64(X & 0x3FFFFFF) // 26 bits
	y := uint64(Y & 0xFFF)     // 12 bits
	z := uint64(Z & 0x3FFFFFF) // 26 bits

	packedPos := ((x << 38) | (y << 26) | z)
	var v [8]byte
	binary.BigEndian.PutUint64(v[:], packedPos)
	return v[:]
}

// https://minecraft.wiki/w/Data_types?oldid=2767030#Position
// 
// decoded as:
//
//	val = read_unsigned_long();
//	x = val >> 38;
//	y = (val >> 26) & 0xFFF;
//	z = val << 38 >> 38;
//
// Note: The details of bit shifting are rather language dependent; the above may work in Java but probably won't in other languages without some tweaking. In particular, you will usually receive positive numbers even if the actual coordinates are negative. This can be fixed by adding something like the following:
//
//	if x >= 2^25 { x -= 2^26 }
//	if y >= 2^11 { y -= 2^12 }
//	if z >= 2^25 { z -= 2^26 }
func DecodePosition(b []byte) (X, Y, Z int32) {
	packedPos := binary.BigEndian.Uint64(b)
	X = int32(packedPos >> 38)
	Y = int32((packedPos >> 26) & 0xFFF)
	Z = int32(packedPos << 38 >> 38) // sign-extend z properly

	if X >= 1<<25 {
		X -= 1 << 26
	}
	if Y >= 1<<11 {
		Y -= 1 << 12
	}
	if Z >= 1<<25 {
		Z -= 1 << 26
	}
	return
}
