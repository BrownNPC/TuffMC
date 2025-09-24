package ds

import "testing"

func TestVarInt(t *testing.T) {
	fortyTwo := EncodeVarInt(42)
	i, _, _ := DecodeVarInt(fortyTwo)
	if i != 42 {
		t.Errorf("Got %v expected 42", i)
	}
}
