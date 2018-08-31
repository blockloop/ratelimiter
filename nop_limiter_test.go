package ratelimiter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNopLimiterReturnsFalse(t *testing.T) {
	l := &NopLimiter{}
	require.False(t, l.Limit("asdf"))
}
