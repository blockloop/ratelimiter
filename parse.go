package ratelimiter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrMalformedLimit is used after attempting to Parse a malformed limit string
	ErrMalformedLimit = errors.New("malformed string")

	// ErrMalformedLimitNumber is used when the first number (the limit) is not a number
	ErrMalformedLimitNumber = errors.New("malformed limit")

	// ErrInvalidLimitNumber is used when the first number (the limit) is < 0
	ErrInvalidLimitNumber = errors.New("limit must be > 0")

	// ErrInvalidDuration is used when the duration is < 0s
	ErrInvalidDuration = errors.New("duration must be > 0s")
)

// ParseLimit parses a limiter string
// limit should be in the format of <count>/<duration>/g where count is the
// maximum count of requests, duration is the length of time, and /g means that
// the limiter is global and not ip specific
//
// Example:
//    1/1m    = one request per minute per IP address
//    10/24h/g = ten requests per day globally
func ParseLimit(limit string) (Limit, error) {
	sp := strings.Split(limit, "/")
	if len(sp) < 2 {
		return Limit{}, ErrMalformedLimit
	}

	max, err := strconv.Atoi(sp[0])
	if err != nil {
		return Limit{}, fmt.Errorf("%s: %s", ErrMalformedLimitNumber, err)
	}
	if max < 1 {
		return Limit{}, ErrInvalidLimitNumber
	}

	dur, err := time.ParseDuration(sp[1])
	if err != nil {
		return Limit{}, fmt.Errorf("%s: %s", "invalid duration", err)
	}
	if dur.Seconds() < 1 {
		return Limit{}, ErrInvalidDuration
	}

	g := len(sp) == 3 && sp[2] == "g"

	return Limit{
		Global: g,
		Limit:  max,
		Dur:    dur,
	}, nil
}

// ParseLimits parses a slice of limits with ParseLimit
func ParseLimits(limits []string) ([]Limit, error) {
	res := make([]Limit, len(limits))
	var err error
	for i, rl := range limits {
		res[i], err = ParseLimit(rl)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

// MustParseLimit calls ParseLimit and panics if there is an error
func MustParseLimit(limit string) Limit {
	l, err := ParseLimit(limit)
	if err != nil {
		panic(err)
	}
	return l
}

// MustParseLimits parses limits and then panics if there is an error
func MustParseLimits(limits []string) []Limit {
	l, err := ParseLimits(limits)
	if err != nil {
		panic(err)
	}
	return l
}
