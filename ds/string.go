package ds

import (
	"fmt"
)

// Decode UTF-8 string prefixed with its size in bytes as a VarInt.
// Maximum length of n characters, which varies by context.
func ReadString(b []byte) (string, int, error) {
	strLen, n, err := ReadVarInt(b)
	if err != nil {
		return "", n, err
	}
	fmt.Println(strLen)
	str := b[n : n+strLen]
	return string(str), n + strLen, nil
}
