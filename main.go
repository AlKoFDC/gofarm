package main

import (
	"net/http"

	"github.com/AlKoFDC/gofarm/gopher"
	"github.com/AlKoFDC/gofarm/server"
)

func main(){
	http.ListenAndServe(":7666", server.Handler(gopher.NewFarm()))
}

// TODO concurrent growing gophers
// TODO request context/cancellation
// TODO gopher kill via context
// TODO prepare actual presentation