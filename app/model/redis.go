package model

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	myerror "moetify/app/error"

	"github.com/go-redis/redis"
	"github.com/speps/go-hashids"
)

// RedisCli is a redis client
type RedisCli struct {
	Cli *redis.Client
}

// URLDetail contains the detail of the shortlink
type URLDetail struct {
	URL                 string        `json:"URL"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

type DetailInfo struct {
	Short               string        `json:"short"`
	Full                string        `json:"full"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

// NewRedisCli creates a new redis client
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
func (r *RedisCli) Shorten(url string, exp int64) (string, error) {
	// convert url to sha1 hash
	h := toHash(url)

	// Only retrive the short URL. No need full or detail
	hres, err := r.Cli.HGet(url, "short").Result()
	if err == redis.Nil {
		// nop
	} else if err != nil {
		log.Fatalln(err)
	} else {
		return hres, nil
	}

	// Create a expiration time variable so that SETEX wouldn't cause millisecond delay from individual SETs
	expireAt := time.Duration(exp)
	mapcontent := map[string]interface{}{
		"short":                 h,
		"full":                  url,
		"created_at":            time.Now().String(),
		"expiration_in_minutes": expireAt,
	}

	batch := func(tx *redis.Tx) error {
		// Queue HMSETs
		err = tx.HMSet(h, mapcontent).Err()
		if err != nil {
			return err
		}
		err = tx.Expire(h, expireAt).Err()
		if err != nil {
			return err
		}

		err = tx.HMSet(url, mapcontent).Err()
		if err != nil {
			return err
		}
		err = tx.Expire(url, expireAt).Err()
		if err != nil {
			return err
		}

		// Exec
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			// Dunno what to do here.
			return nil
		})

		return err
	}

	err = r.Cli.Watch(batch, h, url)
	if err != redis.TxFailedErr {
		return "", err
	}
	return h, nil
}

// ShortLinkInfo gets ShortlinkDetailKey from redis
func (r *RedisCli) ShortLinkInfo(eid string) (*DetailInfo, error) {
	hres, err := r.Cli.HGetAll(eid).Result()
	if err == redis.Nil {
		return nil, myerror.StatusError{
			Code: 404,
			Err:  errors.New("unknown short URL"),
		}
	} else if err != nil {
		return nil, err
	}

	var sli DetailInfo

	tmp, err := json.Marshal(hres)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(tmp, &sli)
	if err != nil {
		return nil, err
	}
	return &sli, nil
}

// Unshorten gets ShortlinkKey from redis
// func (r *RedisCli) Unshorten(eid string) (string, error) {

// 	hres, err := r.Cli.HGet(eid).Result()
// 	if err == redis.Nil {
// 		return "", myerror.StatusError{Code: 404, Err: errors.New("unknown URL")}
// 	} else if err != nil {

// 	}

// 	url, err := r.Cli.Get(eid).Result()
// 	if err == redis.Nil {
// 		return "", myerror.StatusError{Code: 404, Err: errors.New("unknown URL")}
// 	} else if err != nil {
// 		return "", err
// 	} else {
// 		return url, nil
// 	}
// }
