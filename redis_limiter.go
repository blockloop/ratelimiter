package ratelimiter

import (
	"fmt"
	"sort"

	"github.com/garyburd/redigo/redis"
)

var limiter = redis.NewScript(1, `
local limit = tonumber(ARGV[1])
local ttl = tonumber(ARGV[2])
local k = "requests:" .. KEYS[1] .. ":" .. ARGV[2]
local current = redis.call("LLEN", k)
if current >= limit then
    return 1
else
    if redis.call("EXISTS", k) == 1 then
	redis.call("RPUSHX", k, 1)
    else
	redis.call("RPUSH", k, 1)
	redis.call("EXPIRE", k, ttl)
    end
end
return 0`)

type redisPool interface {
	Get() redis.Conn
}

// NewRedisLimiter creates a properly initialized RedisLimiter
func NewRedisLimiter(pool redisPool, limits []Limit) *RedisLimiter {
	// limits must be sorted by TTL descending so that smaller limits don't
	// short circuit the longer ones
	sort.Sort(byDuration(limits))

	return &RedisLimiter{
		Pool:         pool,
		Limits:       limits,
		LimitOnError: true,
	}
}

// RedisLimiter is a rate limit which can evaluate an IP address to determine if it
// should be rate limited using Redis as a backend
type RedisLimiter struct {
	Pool         redisPool
	Limits       []Limit
	LimitOnError bool
	OnError      func(ip string, err error)
}

// Limit checks an IP address to see if it should be ratelimited. It returns
// true if the IP address should be ratelimited and false otherwise any errors
// encountered will return *RedisLimiter.LimitOnError plus the error
func (l *RedisLimiter) Limit(ip string) bool {
	con := l.Pool.Get()
	defer con.Close()

	for _, limit := range l.Limits {
		key := ip
		if limit.Global {
			key = "global"
		}

		limited, err := redis.Bool(limiter.Do(con, key, limit.Limit, limit.Dur.Seconds()))
		if err != nil {
			err := fmt.Errorf("%s: %s", "failed to execute script", err)
			if l.OnError != nil {
				l.OnError(ip, err)
			}

			return l.LimitOnError
		}
		if limited {
			return true
		}
	}
	return false
}
