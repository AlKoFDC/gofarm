package rand

import "math/rand"

func Name() string {
	return randomNames[rand.Intn(len(randomNames)-1)]
}
