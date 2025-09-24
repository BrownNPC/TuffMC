package ds

import (
	"errors"
	"slices"
)

// Decode UTF-8 string prefixed with its size in bytes as a VarInt.
// Maximum length of n characters, which varies by context.
// https://minecraft.wiki/w/Java_Edition_protocol/Data_types#Type:String
func DecodeString(b []byte) (string, int, error) {
	strLen, n, err := DecodeVarInt(b)
	if err != nil {
		return "", n, err
	}
	if n+strLen > len(b) {
		return "", 0, errors.New("string length too long")
	}
	str := b[n : n+strLen]
	return string(str), n + strLen, nil
}
//https://minecraft.wiki/w/Java_Edition_protocol/Data_types#Type:String
func EncodeString(s string) []byte {
	strLen := len(s)
	encodedStrlen := EncodeVarInt(uint(strLen))
	return slices.Concat(encodedStrlen, []byte(s))
}
