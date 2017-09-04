package main

import (
	"context"
	"strings"
	"time"
)

type farm struct{
	gophers map[string]Gopher
	context.Context

	// Request channels
	gopherList chan gopherListRequest
	add chan addRequest
}

type Gopher struct {
	Weight float32
	Height float32
}

func NewGopherFarm(ctx context.Context) farm {
	f := farm{gophers: make(map[string]Gopher), Context: ctx}
	go f.farmer()
	f.gopherList = make(chan gopherListRequest)
	f.add = make(chan addRequest)
	return f
}

func (f farm) Close() {
	close(f.gopherList)
}

type gopherListRequest struct {
	response chan string
}

type addRequest struct {
	name string
	response chan struct{}
}

func (f *farm) farmer() {
	defer f.Close()
	for {
		select {
		case req := <-f.gopherList:
			list := f.synchGopherList()
			req.response<-list
		case req := <-f.add:
			f.synchAdd(req.name)
			req.response<-struct{}{}
		case <-f.Done():
			return
		}
	}
}

func (f farm) Empty() bool {
	return len(f.gophers) <= 0
}

func (f farm) GopherList() string {
	response := make(chan string)
	defer close(response)
	f.gopherList <- gopherListRequest{response: response}
	select {
	case list :=  <-response:
		return list
	case <- f.Done():
		return ""
	}
}
func (f farm) synchGopherList() string {
	if f.Empty() {
		return "empty"
	}
	var list []string
	for name, _ := range f.gophers {
		list = append(list, name)
	}
	return "- " + strings.Join(list, "\n- ")
}

func (f *farm) Add(name string) {
	response := make(chan struct{})
	defer close(response)
	f.add <- addRequest{name: name, response: response}
	select {
	case <-response:
	case <-f.Done():
	}
	return
	
}
func (f *farm) synchAdd(name string) {
	if _, ok := f.gophers[name]; ok {
		return
	}
	hour, min, sec := time.Now().Clock()
	f.gophers[name] = Gopher{
		Height: float32(sec)/float32(hour),
		Weight: float32(sec)/float32(min),
	}
}
