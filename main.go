package main

import (
	"net/http"

	"github.com/AlKoFDC/gofarm/gopher"
	"github.com/AlKoFDC/gofarm/server"
)

func main() {
	http.ListenAndServe(":7666", server.Handler(gopher.NewFarm()))
}

var _ server.Farmer = (*gopher.Farm)(nil)
