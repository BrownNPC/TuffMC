package ds

func EncodeInt(v int32) []byte {
	return []byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}
func DecodeInt32(b []byte) int32 {
	if len(b) < 4 {
		panic("not enough bytes to decode int32")
	}
	return int32(b[0])<<24 |
		int32(b[1])<<16 |
		int32(b[2])<<8 |
		int32(b[3])
}
