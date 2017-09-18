package gopher

import "fmt"

const lengthName = 12

type gopher struct {
	name string
	skill int
	rate int
}

func (g gopher) String() string {
	return fmt.Sprintf("%s %d %d", g.name, g.skill, g.rate)
}

func spawn() gopher {
	return gopher{
		name: randomGopherName(lengthName),
	}
}

