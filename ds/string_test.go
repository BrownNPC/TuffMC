package ds

import "testing"

func TestString(t *testing.T) {
	str := "The cake is a lie"
	encodedStr := WriteString(str)

	lie, _, _ := ReadString(encodedStr)
	if lie != str {
		t.Error("Strings did not match")
	}
}
