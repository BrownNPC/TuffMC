package ds

import "testing"

func TestString(t *testing.T) {
	str := "The cake is a lie"
	encodedStr := EncodeString(str)

	lie, _, _ := DecodeString(encodedStr)
	if lie != str {
		t.Error("Strings did not match")
	}
}
