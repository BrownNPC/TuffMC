package ds

import (
	"bufio"
	"errors"
	"fmt"
)

// https://minecraft.wiki/w/Minecraft_Wiki:Projects/wiki.vg_merge/Protocol?oldid=2992295#VarInt_and_VarLong
func EncodeVarInt(value uint) []byte {
	const SEGMENT_BITS = 0x7F
	const CONTINUE_BIT = 0x80
	var buf []byte
	for {
		if value&^SEGMENT_BITS == 0 {
			buf = append(buf, byte(value))
			return buf
		}

		buf = append(buf,
			byte(value)&SEGMENT_BITS|CONTINUE_BIT,
		)
		value >>= 7
	}
}

// https://minecraft.wiki/w/Minecraft_Wiki:Projects/wiki.vg_merge/Protocol?oldid=2992295#VarInt_and_VarLong
func DecodeVarInt(b []byte) (value, n int, err error) {
	const SEGMENT_BITS = 0x7F
	const CONTINUE_BIT = 0x80

	var position int
	for i, currentByte := range b {

		value |= int(currentByte&SEGMENT_BITS) << position

		if currentByte&CONTINUE_BIT == 0 {
			return value, i + 1, err
		}

		position += 7
		if position >= 32 {
			return 0, 0, errors.New("VarInt is too big")
		}
	}
	err = errors.New("buffer too small")
	return
}

// https://minecraft.wiki/w/Minecraft_Wiki:Projects/wiki.vg_merge/Protocol?oldid=2992295#VarInt_and_VarLong
func DecodeVarIntFromReader(r *bufio.Reader) (v int, err error) {
	const SEGMENT_BITS = 0x7F
	const CONTINUE_BIT = 0x80

	var position, value uint
	var currentByte byte
	for {
		currentByte, err = r.ReadByte()
		if err != nil {
			return
		}
		value |= (uint(currentByte) & SEGMENT_BITS) << position
		if currentByte&CONTINUE_BIT == 0 {
			break
		}
		position += 7
		if position > 32 {
			return 0, fmt.Errorf("varint too big")
		}
	}
	return int(value), nil
}
