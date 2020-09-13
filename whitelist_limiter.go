package ratelimiter

import (
	"fmt"
	"net"
)

// WhitelistedLimiter wraps a ratelimiter with a whitelist
type WhitelistedLimiter struct {
	Limiter     Limiter
	Whitelist   []*net.IPNet
	OnWhitelist func(ip string)
}

// ParseWhitelist parses a list of strings as CIDRs
func ParseWhitelist(whitelist []string) ([]*net.IPNet, error) {
	wl := make([]*net.IPNet, 0, len(whitelist))
	for _, cidr := range whitelist {
		_, c, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CIDR %q: %w", cidr, err)
		}
		wl = append(wl, c)
	}
	return wl, nil
}

// NewWhitelistedLimiter constructs a new WhitelistedLimiter.
func NewWhitelistedLimiter(limiter Limiter, whitelist []*net.IPNet) *WhitelistedLimiter {
	return &WhitelistedLimiter{
		Limiter:   limiter,
		Whitelist: whitelist,
	}
}

// Limit checks the whitelist for whitelisted IP addresses and then return
// false if any match. If none match then it defers to w.Limiter.Limit
func (w *WhitelistedLimiter) Limit(ip string) bool {
	for _, wl := range w.Whitelist {
		if wl.Contains(net.ParseIP(ip)) {
			if w.OnWhitelist != nil {
				w.OnWhitelist(ip)
			}
			return false
		}
	}

	return w.Limiter.Limit(ip)
}
