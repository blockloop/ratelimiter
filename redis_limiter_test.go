package ratelimiter

import (
	"math/rand"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/garyburd/redigo/redis"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestLimiter(t *testing.T) {
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	for i := 0; i < 10; i++ {
		l := Limit{
			Dur:    time.Duration(rand.Intn(10)+1) * time.Second,
			Global: rand.Intn(10)%2 == 0,
			Limit:  rand.Intn(10) + 1,
		}

		t.Run(l.String(), limitTestFunc(srv, l))
		srv.FlushAll()
	}

}

func limitTestFunc(r *miniredis.Miniredis, l Limit) func(t *testing.T) {
	return func(t *testing.T) {
		limiter := NewRedisLimiter(&fakePool{addr: r.Addr()}, []Limit{l})
		ip := uuid.New()

		for i := 0; i < l.Limit; i++ {
			assert.False(t, limiter.Limit(ip), l.String())
		}

		assert.True(t, limiter.Limit(ip), l.String())
		r.FastForward(l.Dur)
		assert.False(t, limiter.Limit(ip))
	}
}

func TestLimiterWhitelistNeverLimits(t *testing.T) {
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	l := Limit{
		Dur:    time.Minute,
		Global: true,
		Limit:  1,
	}

	p := &fakePool{addr: srv.Addr()}
	ip := uuid.New()
	limiter := NewRedisLimiter(p, []Limit{l}, ip)

	assert.False(t, limiter.Limit(ip), l.String())
	assert.False(t, limiter.Limit(ip), l.String())
}

func TestLimiterWhitelistCallsOnWhitelist(t *testing.T) {
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	l := Limit{
		Dur:    time.Minute,
		Global: true,
		Limit:  1,
	}

	var result string

	p := &fakePool{addr: srv.Addr()}
	ip := uuid.New()
	limiter := NewRedisLimiter(p, []Limit{l}, ip)
	limiter.OnWhitelist = func(ip string) {
		result = ip
	}

	assert.False(t, limiter.Limit(ip), l.String())
	assert.Equal(t, ip, result)
}

func TestLimiterCallsOnError(t *testing.T) {
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	l := Limit{
		Dur:    time.Minute,
		Global: true,
		Limit:  1,
	}

	var result string

	p := &deadPool{addr: srv.Addr()}
	ip := uuid.New()
	limiter := NewRedisLimiter(p, []Limit{l})
	limiter.OnError = func(ip string, err error) {
		assert.Error(t, err)
		result = ip
	}

	assert.True(t, limiter.Limit(ip), l.String())
	assert.Equal(t, ip, result)
}

func TestLimiterRespectsLimitOnError(t *testing.T) {
	srv, err := miniredis.Run()
	require.NoError(t, err)
	defer srv.Close()

	l := Limit{
		Dur:    time.Minute,
		Global: true,
		Limit:  1,
	}

	var result string

	p := &deadPool{addr: srv.Addr()}
	ip := uuid.New()

	limiter := NewRedisLimiter(p, []Limit{l})
	limiter.OnError = func(ip string, err error) {
		assert.Error(t, err)
		result = ip
	}

	t.Run("true", func(t *testing.T) {
		limiter.LimitOnError = true
		assert.True(t, limiter.Limit(ip), l.String())
		assert.Equal(t, ip, result)
	})

	t.Run("false", func(t *testing.T) {
		limiter.LimitOnError = false
		assert.False(t, limiter.Limit(ip), l.String())
		assert.Equal(t, ip, result)
	})
}

type fakePool struct {
	addr string
}

func (f *fakePool) Get() redis.Conn {
	con, err := redis.DialURL("redis://" + f.addr)
	if err != nil {
		panic(err)
	}

	return con
}

type deadPool struct {
	addr string
}

func (f *deadPool) Get() redis.Conn {
	con, err := redis.DialURL("redis://" + f.addr)
	if err != nil {
		panic(err)
	}
	// close the error to kill the connection
	con.Close()

	return con
}
