package ratelimiter

// WhitelistedLimiter wraps a ratelimiter with a whitelist
type WhitelistedLimiter struct {
	Limiter     Limiter
	Whitelist   []string
	OnWhitelist func(ip string)
}

// NewWhitelistedLimiter constructs a new WhitelistedLimiter
func NewWhitelistedLimiter(limiter Limiter, whitelist []string) *WhitelistedLimiter {
	return &WhitelistedLimiter{
		Limiter:   limiter,
		Whitelist: whitelist,
	}
}

// Limit checks the whitelist for whitelisted IP addresses and then return
// false if any match. If none match then it defers to w.Limiter.Limit
func (w *WhitelistedLimiter) Limit(ip string) bool {
	for _, wl := range w.Whitelist {
		if ip == wl {
			if w.OnWhitelist != nil {
				w.OnWhitelist(ip)
			}
			return false
		}
	}

	return w.Limiter.Limit(ip)
}
