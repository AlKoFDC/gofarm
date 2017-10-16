package main

import (
	"github.com/AlKoFDC/gofarm/gopher"
	"github.com/AlKoFDC/gofarm/server"
)

var _ server.Farmer = (*gopher.Farm)(nil)
