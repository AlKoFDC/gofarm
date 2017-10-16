package gopher

import (
	"testing"

	"context"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSpawnGophers(t *testing.T) {
	newGopher1 := spawn(context.Background())
	newGopher2 := spawn(context.Background())
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

func TestGrow(t *testing.T) {
	babyGopher := spawn(context.Background())
	initSkill := babyGopher.skill
	time.Sleep(sleep * 2)
	if gotSkill := babyGopher.skill; initSkill >= gotSkill {
		t.Errorf("baby gopher failed to learn a damn thing!")
	}

	babyGopher.skill = maxSkill
	time.Sleep(sleep * 2)
	if gotSkill := babyGopher.skill; gotSkill > maxSkill {
		t.Errorf("baby gopher learned too damn much!")
	}
}
