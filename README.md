# Go Rate Limiter

[![GoDoc](https://godoc.org/github.com/blockloop/ratelimiter?status.svg)](https://godoc.org/github.com/blockloop/ratelimiter)
[![Travis](https://img.shields.io/travis/blockloop/ratelimiter.svg)](https://travis-ci.org/blockloop/ratelimiter)
[![Coveralls github](https://img.shields.io/coveralls/github/blockloop/ratelimiter.svg)](https://coveralls.io/github/blockloop/ratelimiter)
[![Report Card](https://goreportcard.com/badge/github.com/blockloop/ratelimiter)](https://goreportcard.com/report/github.com/blockloop/ratelimiter)

`ratelimiter` is a multi-ratelimiter for go. It is easily configurable, allows whitelist, and easy to use.

## Rate Limiters

```go
type Limiter interface {
        Limit(string) bool
}
```

You create a limiter by providing it with a configured Limit which is either manually created or parsed from a string

```go
// Using a parsed limit
rl := ratelimiter.MustParseLimit("10/1m")
limiter := ratelimiter.NewRedisLimiter(redisPool, rls)

// Manually created limit
// one per minute, non-global (IP specific)
limit := ratelimiter.Limit{
        Global: false,
        Limit: 1,
        Dur: time.Minute
}
limiter := ratelimiter.NewRedisLimiter(redisPool, limit)

```

## Limits

Limits are defined as the following struct

```go
type Limit struct {
	Global bool
	Limit  int
	Dur    time.Duration
}
```

You can either manually create them or you can use the parser. The parser parses strings separated into three sections

- Section 1: How many requests to limit to
- Section 2: For what duration
- Section 3: (optional) is this limit global?

### Example limit strings

```go
// create a Limit which limits requests to 1 per second by identifier (IP address)
rl := ratelimiter.MustParseLimit("1/1s")

// create a Limit which limits requests to 10 per 5 minutes by identifier (IP address)
rl := ratelimiter.MustParseLimit("10/5m")

// create a Limit which limits requests to 1000 per day globally
rl := ratelimiter.MustParseLimit("1000/24h/g")

// create chained Limits
rls := ratelimiter.MustParseLimits([]string{"1000/24h/g", "1/1s", "10/1m"})
```

Ratelimits are executed in descending order of duration. If you configure the following two limits:

1. 1/1s
2. 10/1h

They will be executed in reverse order because the larger groups should be checked first. Global rate limits are also executed first.


## Whitelist

To whitelist IP addresses use the `WhitelistedLimiter` which wraps a real limiter

```go
rl := ratelimiter.MustParseLimit("10/1m")
limiter := ratelimiter.NewRedisLimiter(redisPool, rls)

whitelist := []string{"127.0.0.1"}
wl := ratelimiter.NewWhitelistedLimiter(limiter, whitelist)
```

## Implemented Limiters

For now only `RedisRatelimiter` is implemented, however more can be added. Feel free to submit a PR.
