package ratelimiter

import (
	"testing"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWhitelistedLimiterIgnoresLimiterWhenWhitelisted(t *testing.T) {
	ip := uuid.New()
	fake := &fakeLimiter{
		LimitFunc: func(string) bool {
			return true
		},
	}

	wl := NewWhitelistedLimiter(fake, []string{ip})

	assert.False(t, wl.Limit(ip))
}

func TestWhitelistedLimiterDefersWhenNotWhitelisted(t *testing.T) {
	ip := uuid.New()
	done := make(chan struct{})
	fake := &fakeLimiter{
		LimitFunc: func(in string) bool {
			defer close(done)
			assert.Equal(t, ip, in)
			return true
		},
	}

	wl := NewWhitelistedLimiter(fake, nil)

	assert.True(t, wl.Limit(ip))
}

func TestWhitelistedLimiterCallsOnWhitelistWhenWhitelisted(t *testing.T) {
	ip := uuid.New()
	done := make(chan struct{})
	fake := &fakeLimiter{}

	wl := NewWhitelistedLimiter(fake, []string{ip})
	wl.OnWhitelist = func(in string) {
		defer close(done)
		assert.Equal(t, ip, in)
	}

	assert.False(t, wl.Limit(ip))
}
