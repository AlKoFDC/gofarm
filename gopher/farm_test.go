package gopher

import (
	"context"
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestCountGophers(t *testing.T) {
	populatedFarm := NewFarm()
	populatedFarm.stable = map[int]*gopher{
		1: {},
		2: {},
	}

	for _, test := range []struct {
		desc   string
		farm   *Farm
		expect int
	}{
		{desc: "empty", farm: NewFarm(), expect: 0},
		{desc: "two gophers", farm: populatedFarm, expect: 2},
	} {
		t.Run(test.desc, func(t *testing.T) {
			if count := len(test.farm.stable); count != test.expect {
				t.Errorf("expected %d but got %d", test.expect, count)
			}
		})
	}
}

func TestAddGophers(t *testing.T) {
	testFarm := NewFarm()
	testFarm.Add(context.Background())
	if count := len(testFarm.stable); count != 1 {
		t.Errorf("expected first gopher to be added to the farm, but got total %d", count)
	}
	testFarm.Add(context.Background())
	if count := len(testFarm.stable); count != 2 {
		t.Errorf("expected second gopher to be added to the farm, but got total %d", count)
	}
}

func TestList(t *testing.T) {
	testFarm := NewFarm()
	testFarm.Add(context.Background())
	testFarm.Add(context.Background())
	testFarm.Add(context.Background())
	count := len(testFarm.stable) + 2 // plus two title lines
	if list := testFarm.List(); len(list) != count {
		t.Errorf("expected to get the list of %d gophers, but got %d", count, len(list))
	}
}

func TestKill(t *testing.T) {
	populatedFarm := NewFarm()
	populatedFarm.add(context.Background())
	populatedFarm.add(context.Background())

	const (
		ghostGopher = math.MaxInt64
		realGopher  = 0
	)

	err := populatedFarm.Kill(ghostGopher)
	if err == nil {
		t.Errorf("expected to get an error trying to kill the ghost gopher id %d, but got none", ghostGopher)
	}
	if !strings.Contains(err.Error(), strconv.Itoa(ghostGopher)) {
		t.Errorf("expected to get an error mentioning the ghost gopher id %d, but got %s", ghostGopher, err)
	}

	err = populatedFarm.Kill(realGopher)
	if err != nil {
		t.Errorf("expected to get no error trying to kill the real gopher id %d, but got %s", realGopher, err)
	}

	if expCount, c := 1, len(populatedFarm.stable); c != expCount {
		t.Errorf("expected number of gophers alive to be %d, but got %d", expCount, c)
	}
}
