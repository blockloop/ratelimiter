package ratelimiter

import (
	"fmt"
	"time"
)

// Limiter checks an ip address to see if it should be ratelimited
type Limiter interface {
	Limit(ip string) bool
}

// Limit is a limiter used with New to execuate a ratelimiter
type Limit struct {
	Global bool
	Limit  int
	Dur    time.Duration
}

func (l *Limit) String() string {
	s := fmt.Sprintf("%d/%s", l.Limit, l.Dur)
	if l.Global {
		return s + "/g"
	}

	return s
}

type ByDuration []Limit

func (d ByDuration) Len() int      { return len(d) }
func (d ByDuration) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d ByDuration) Less(i, j int) bool {
	a, b := d[i], d[j]
	if a.Global != b.Global {
		// place globals at the end so they aren't affected by local limits
		return b.Global
	}

	return a.Dur > b.Dur
}
