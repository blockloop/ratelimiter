package ratelimiter

// NopLimiter is a limiter that returns false for all Limit() queries
type NopLimiter struct{}

// Limit always returns false
func (l *NopLimiter) Limit(ip string) bool {
	return false
}
