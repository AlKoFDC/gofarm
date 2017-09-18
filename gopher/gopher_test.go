package gopher

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSpawnGophers(t *testing.T) {
	newGopher1 := spawn()
	newGopher2 := spawn()
	if cmp.Diff(newGopher1, newGopher2, cmp.AllowUnexported(gopher{})) == "" {
		t.Error("expected to get two different gophers, but got the same")
	}
	emptyGopher := gopher{}
	if cmp.Diff(newGopher1, emptyGopher, cmp.AllowUnexported(gopher{})) == "" {
		t.Error("expected to get proper gopher #1, but got empty one")
	}
	if cmp.Diff(newGopher2, emptyGopher, cmp.AllowUnexported(gopher{})) == "" {
		t.Error("expected to get proper gopher #2, but got empty one")
	}
}

