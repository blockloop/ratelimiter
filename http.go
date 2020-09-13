package ratelimiter

import (
	"net/http"
	"strings"
)

// Middleware creates a new rate limiter for HTTP
func Middleware(l Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if l.Limit(strings.Split(r.RemoteAddr, ":")[0]) {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
