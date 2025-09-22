package ds

import (
	"errors"
	"slices"
)

// Decode UTF-8 string prefixed with its size in bytes as a VarInt.
// Maximum length of n characters, which varies by context.
func ReadString(b []byte) (string, int, error) {
	strLen, n, err := ReadVarInt(b)
	if err != nil {
		return "", n, err
	}
	if n+strLen > len(b[n:]) {
		return "", 0, errors.New("string length too long")
	}
	str := b[n : n+strLen]
	return string(str), n + strLen, nil
}
func WriteString(s string) []byte {
	strLen := len(s)
	encodedStrlen := WriteVarInt(uint(strLen))
	return slices.Concat(encodedStrlen, []byte(s))
}
