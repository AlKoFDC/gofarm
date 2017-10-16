package rand

import "testing"

func TestName(t *testing.T){
	length := 13
	newString := String(length)
	if len(newString) != length {
		t.Errorf("failed to get random string of requested length %d, got: %s", length, newString)
	}
}
