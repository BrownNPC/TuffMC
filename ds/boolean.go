package ds

// Either false or true	True is encoded as 0x01, false as 0x00.
func EncodeBool(v bool) []byte {
	if v == true {
		return []byte{1}
	} else {
		return []byte{0}
	}
}
func DecodeBool(b byte) bool {
	if b == 1 {
		return true
	}
	if b == 0 {
		return false
	}
	panic("boolean must be 0 or 1")
}
