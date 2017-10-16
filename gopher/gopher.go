package gopher

import (
	"context"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/AlKoFDC/gofarm/rand"
	"github.com/AlKoFDC/gofarm/calltreeid"
)

const (
	lengthName     = 12
	maxSkill       = 999
	minSkillInit   = 100
	maxSkillRange  = 100
	maxRate        = 10
	sleep          = timeSpeed * time.Millisecond
	statusLearning = "LEARNING"
	statusRetired  = "RETIRED"
)

type gopher struct {
	name      string
	skill     int
	rate      int
	terminate func()
	status    string
}

func (g gopher) String() string {
	return fmt.Sprintf("%20s   %d    %d %s", g.name, g.skill, g.rate, g.status)
}

func spawn(ctx context.Context) *gopher {
	gopherContext, cancel := context.WithCancel(context.Background())

	//cti := calltreeid.FromContext(ctx)
	//gopherContext = calltreeid.InNewContext(gopherContext, cti)

	babyGopher := &gopher{
		name:      rand.Name(),
		skill:     minSkillInit + mrand.Intn(maxSkillRange),
		rate:      mrand.Intn(maxRate),
		terminate: cancel,
	}
	babyGopher.grow(gopherContext, sleep)
	return babyGopher
}

func (g *gopher) grow(ctx context.Context, sleep time.Duration) {
	g.status = statusLearning
	go func() {
		for {
			select {
			case <-ctx.Done():
				g.status = statusRetired
				return
			case <-time.After(sleep):
				g.skill += g.rate
				if g.skill > maxSkill {
					g.skill = maxSkill
				}
			}
		}
	}()
}

func (g *gopher) slaughter() {
	g.terminate()
}
