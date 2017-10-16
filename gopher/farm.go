package gopher

import (
	"context"
	"fmt"
	"sort"
)

const timeSpeed = 500

type Farm struct {
	stable    map[int]*gopher
	actions   chan func()
	nextIndex int
}

func NewFarm() *Farm {
	farm := &Farm{
		stable:  make(map[int]*gopher),
		actions: make(chan func()),
	}
	go func() {
		for action := range farm.actions {
			action()
		}
	}()
	return farm
}

func (f Farm) List() []string {
	response := make(chan []string)
	f.actions <- func() {
		response <- f.list()
	}
	return <-response
}

func (f *Farm) Kill(index int) error {
	response := make(chan error)
	f.actions <- func() {
		response <- f.kill(index)
	}
	return <-response
}

func (f *Farm) Add(ctx context.Context) string {
	response := make(chan string)
	f.actions <- func() {
		response <- f.add(ctx)
	}
	return <-response
}

func (f Farm) list() []string {
	list := make([]string, 0, len(f.stable))
	for id, g := range f.stable {
		list = append(list, fmt.Sprintf("%2d %s", id, g.String()))
	}
	sort.Strings(list)
	if len(list) < 1 {
		return []string{}
	}
	return append(title, list...)
}

func (f *Farm) kill(index int) error {
	gopher, ok := f.stable[index]
	if !ok {
		return fmt.Errorf("could not find a gopher with id %d to kill", index)
	}
	gopher.slaughter()
	delete(f.stable, index)
	return nil
}

func (f *Farm) add(ctx context.Context) string {
	newGopher := spawn(ctx)
	f.stable[f.nextIndex] = newGopher
	f.nextIndex++
	return newGopher.name
}

var title = []string{
	fmt.Sprintf("%s %20s %s %s %s", "ID", "GOPHER NAME", "SKILL", "RATE", "STATUS"),
	"----------------------------------------------------",
}
