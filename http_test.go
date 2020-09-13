package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddlewareCallsNextWhenLimitIsFalse(t *testing.T) {
	const code = 201
	next := &fakeHandler{
		ServeHTTPFunc: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
		},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:22826"

	mw := Middleware(&fakeLimiter{})

	mw(next).ServeHTTP(w, r)
	w.Flush()
	require.EqualValues(t, code, w.Code)
}

func TestMiddlewarePassesOnlyIPFromRemoteAddrToLimiter(t *testing.T) {
	next := &fakeHandler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:22826"

	done := make(chan struct{})
	limiter := &fakeLimiter{
		LimitFunc: func(ip string) bool {
			defer close(done)
			require.EqualValues(t, strings.Split(r.RemoteAddr, ":")[0], ip)
			return true
		},
	}

	mw := Middleware(limiter)

	mw(next).ServeHTTP(w, r)
	w.Flush()
	<-done
	require.EqualValues(t, http.StatusTooManyRequests, w.Code)
}

func TestMiddlewareSetsStatusCodeToTooManyRequestsWhenLimiterReturnsTrue(t *testing.T) {
	next := &fakeHandler{}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "127.0.0.1:22826"

	limiter := &fakeLimiter{
		LimitFunc: func(string) bool {
			return true
		},
	}

	mw := Middleware(limiter)

	mw(next).ServeHTTP(w, r)
	w.Flush()
	require.EqualValues(t, http.StatusTooManyRequests, w.Code)
}

type fakeHandler struct {
	ServeHTTPFunc func(http.ResponseWriter, *http.Request)
}

func (f *fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.ServeHTTPFunc != nil {
		f.ServeHTTPFunc(w, r)
		return
	}
}

type fakeLimiter struct {
	LimitFunc func(string) bool
}

func (f *fakeLimiter) Limit(ip string) bool {
	if f.LimitFunc != nil {
		return f.LimitFunc(ip)
	}
	return false
}
