package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
	"github.com/speps/go-hashids"
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

func toHash(url string) string {
	hd := hashids.NewData()
	hd.Salt = url
	hd.MinLength = 0
	h, _ := hashids.NewWithData(hd)
	r, _ := h.Encode([]int{45, 434, 1313, 99})
	return r
}

// Shorten convert url to shortlink
func (r RedisCli) Shorten(url string, exp int64) (string, error) {
	// convert url to sha1 hash
	h := toHash(url)

	// fetch it if url is cached
	d, err := r.Cli.Get(fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil {
		// cache doesn't exist, do nothing
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			// expiration, do nothing
		} else {
			return d, nil
		}
	}

	// increase the global counter
	err = r.Cli.Incr(URLIDKEY).Err()
	if err != nil {
		return "", err
	}

	// encode global counter to base62
	id, err := r.Cli.Get(URLIDKEY).Int64()
	if err != nil {
		return "", err
	}
	eid := base62.EncodeInt64(id)

	// store the encoded id
	err = r.Cli.Set(fmt.Sprintf(ShortlinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	// store the hash of url
	err = r.Cli.Set(fmt.Sprintf(URLHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	detail, err := json.Marshal(&URLDetail{
		URL:                 url,
		CreatedAt:           time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})
	if err != nil {
		return "", err
	}

	// store the detail
	err = r.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, eid), detail, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	return eid, nil
}

func (r *RedisCli) ShortLinkInfo(eid string) (interface{}, error) {
	d, err := r.Cli.Get(fmt.Sprintf(ShortlinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{404, errors.New("Unknown short URL")}
	} else if err != nil {
		return "", err
	} else {
		return d, nil
	}
}

func (r *RedisCli) Unshorten(eid string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{404, errors.New("Unknown URL")}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
}
