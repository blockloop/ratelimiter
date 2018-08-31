package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseLimitErrors(t *testing.T) {
	tests := map[string]string{
		"asdf":      "malformed string",
		"1/":        "malformed string",
		"/jj//":     "malformed string",
		"-1/1m":     "negative limit",
		"1/1y":      "invalid duration",
		"1/0s/g":    "zero duration",
		"1/200ms/g": "< 1s duration",
	}

	for test, msg := range tests {
		t.Run(msg, func(t *testing.T) {
			_, err := ParseLimit(test)
			assert.Error(t, err)
		})
	}
}

func TestParseLimit(t *testing.T) {
	tests := map[string]Limit{
		"1/1s": Limit{
			Dur:    time.Second,
			Global: false,
			Limit:  1,
		},
		"10/1m": Limit{
			Dur:    time.Minute,
			Global: false,
			Limit:  10,
		},
		"10/24h/g": Limit{
			Dur:    24 * time.Hour,
			Global: true,
			Limit:  10,
		},
		"1/1s/g": Limit{
			Dur:    time.Second,
			Global: true,
			Limit:  1,
		},
	}

	for raw, limit := range tests {
		t.Run(raw, func(t *testing.T) {
			res, err := ParseLimit(raw)
			assert.NoError(t, err)
			assert.EqualValues(t, limit, res)
		})
	}
}

func TestParseLimits(t *testing.T) {
	raw := []string{"1/1s", "10/1m", "10/24h/g", "1/1s/g"}
	expected := []Limit{
		Limit{
			Dur:    time.Second,
			Global: false,
			Limit:  1,
		},
		Limit{
			Dur:    time.Minute,
			Global: false,
			Limit:  10,
		},
		Limit{
			Dur:    24 * time.Hour,
			Global: true,
			Limit:  10,
		},
		Limit{
			Dur:    time.Second,
			Global: true,
			Limit:  1,
		},
	}

	for i, l := range raw {
		t.Run(l, func(t *testing.T) {
			res, err := ParseLimit(l)
			assert.NoError(t, err)
			assert.EqualValues(t, expected[i], res)
		})
	}
}

func TestMustParseLimitsPanicsWhenError(t *testing.T) {
	assert.Panics(t, func() {
		MustParseLimits([]string{"asdfasdf", "jfjf"})
	})
}

func TestMustParseLimitPanicsWhenError(t *testing.T) {
	assert.Panics(t, func() {
		MustParseLimit("asdfasdf")
	})
}

func TestMustParseLimitsParses(t *testing.T) {
	raw := []string{"1/1s", "10/1m"}
	expected := []Limit{
		Limit{
			Dur:    time.Second,
			Global: false,
			Limit:  1,
		},
		Limit{
			Dur:    time.Minute,
			Global: false,
			Limit:  10,
		},
	}

	actual := MustParseLimits(raw)
	assert.EqualValues(t, expected, actual)
}

func TestMustParseLimitParses(t *testing.T) {
	raw := "1/1s"
	expected := Limit{
		Dur:    time.Second,
		Global: false,
		Limit:  1,
	}

	res := MustParseLimit(raw)
	assert.EqualValues(t, expected, res)
}
