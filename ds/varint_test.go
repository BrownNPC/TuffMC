package ds

import "testing"

func TestVarInt(t *testing.T) {
	fortyTwo := WriteVarInt(42)
	i, _, _ := ReadVarInt(fortyTwo)
	if i != 42 {
		t.Errorf("Got %v expected 42", i)
	}
}
