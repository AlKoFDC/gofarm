package gopher

import "testing"

func TestName(t *testing.T){
	newName := randomGopherName(lengthName)
	if len(newName) != lengthName {
		t.Errorf("failed to get random gopher name of requested length %d, got: %s", lengthName, newName)
	}
}
