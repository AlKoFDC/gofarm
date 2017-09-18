package gopher

import (
	"fmt"
)

type Farm struct {
	gophers   map[int]gopher
	nextIndex int
	actions   chan func()
}

func NewFarm() *Farm {
	f := &Farm{
		gophers: make(map[int]gopher),
		actions: make(chan func()),
	}
	go func() {
		for a := range f.actions {
			a()
		}
	}()
	return f
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

func (f *Farm) Add() string {
	response := make(chan string)
	f.actions <- func() {
		response <- f.add()
	}
	return <-response
}

func (f Farm) list() []string {
	list := make([]string, 0, f.count())
	for id, g := range f.gophers {
		list = append(list, fmt.Sprintf("%d: %s", id, g.String()))
	}
	return list
}

func (f *Farm) kill(index int) error {
	if _, ok := f.gophers[index]; !ok {
		return fmt.Errorf("could not find a gopher with id %d to kill", index)
	}
	delete(f.gophers, index)
	return nil
}

func (f *Farm) add() string {
	newGopher := spawn()
	f.gophers[f.nextIndex] = newGopher
	f.nextIndex++
	return newGopher.name
}

func (f Farm) count() int {
	return len(f.gophers)
}
