package main

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	// URLIDKEY is the global counter
	URLIDKEY = "next.url.id"
	// ShortlinkKey maps the shortlint to the original URL
	ShortlinkKey = "shortlink:%s:url"
	// URLHashKey maps the hash of the URL to the shortlink
	URLHashKey = "urlhashkey:%s:url"
	// ShortlinkDetailKey map the shortlink and URL info
	ShortlinkDetailKey = "shortlink:%s:detail"
)

// RedisClient is a redis client
type RedisCli struct {
	Cli *redis.Client
}

// URLDetail contains the detail of the shortlink
type URLDetail struct {
	URL                 string        `json:"URL"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

// NewRedisClient creates a new redis client
func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisCli{Cli: c}
}
