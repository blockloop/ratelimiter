package ratelimiter

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestByDurationSortsByDurationDescending(t *testing.T) {
	unsorted := []Limit{
		{Dur: time.Minute},
		{Dur: time.Minute},
		{Dur: time.Second},
		{Dur: time.Hour},
	}

	sorted := []Limit{
		{Dur: time.Hour},
		{Dur: time.Minute},
		{Dur: time.Minute},
		{Dur: time.Second},
	}

	sort.Sort(ByDuration(unsorted))
	assert.Equal(t, sorted, unsorted)
}

func TestByDurationSortsGlobalBeforeNonGlobal(t *testing.T) {
	unsorted := []Limit{
		{
			Dur:    time.Minute,
			Global: true,
		},
		{
			Dur:    time.Minute,
			Global: false,
		},
		{
			Dur:    time.Hour,
			Global: false,
		},
		{
			Dur:    time.Minute,
			Global: true,
		},
	}

	sorted := []Limit{
		{
			Dur:    time.Hour,
			Global: false,
		},
		{
			Dur:    time.Minute,
			Global: false,
		},
		{
			Dur:    time.Minute,
			Global: true,
		},
		{
			Dur:    time.Minute,
			Global: true,
		},
	}

	sort.Sort(ByDuration(unsorted))
	assert.Equal(t, sorted, unsorted)
}
