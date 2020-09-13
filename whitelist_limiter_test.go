package ratelimiter

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhitelistedLimiterIgnoresLimiterWhenWhitelisted(t *testing.T) {
	ip, cidr, err := net.ParseCIDR("192.168.1.100/32")
	require.NoError(t, err)
	fake := &fakeLimiter{
		LimitFunc: func(string) bool {
			return true
		},
	}

	wl := NewWhitelistedLimiter(fake, []*net.IPNet{cidr})
	require.NoError(t, err)

	assert.False(t, wl.Limit(ip.String()))
}

func TestWhitelistedLimiterDefersWhenNotWhitelisted(t *testing.T) {
	ip, _, err := net.ParseCIDR("192.168.1.100/32")
	require.NoError(t, err)
	done := make(chan struct{})
	fake := &fakeLimiter{
		LimitFunc: func(in string) bool {
			defer close(done)
			assert.Equal(t, ip.String(), in)
			return true
		},
	}

	wl := NewWhitelistedLimiter(fake, nil)
	require.NoError(t, err)

	assert.True(t, wl.Limit(ip.String()))
}

func TestWhitelistedLimiterCallsOnWhitelistWhenWhitelisted(t *testing.T) {
	ip, cidr, err := net.ParseCIDR("192.168.1.100/32")
	require.NoError(t, err)
	done := make(chan struct{})
	fake := &fakeLimiter{}

	wl := NewWhitelistedLimiter(fake, []*net.IPNet{cidr})
	require.NoError(t, err)
	wl.OnWhitelist = func(in string) {
		defer close(done)
		assert.Equal(t, ip.String(), in)
	}

	assert.False(t, wl.Limit(ip.String()))
}
