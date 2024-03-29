package model

import (
	"errors"
	"log"
	"time"

	myerror "github.com/kakugirai/moetify/app/error"

	"github.com/go-redis/redis"
	"github.com/speps/go-hashids"
)

// RedisCli is a redis client
type RedisCli struct {
	Cli *redis.Client
}

// DetailInfo contains the detail of the shortlink
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
	//log.Println("hashed", h)

	// Only retrieve the short URL. No need full or detail
	hres, err := r.Cli.HGet(url, "short").Result()
	//log.Println("hget", hres, err)
	if err == redis.Nil {
		// nop
	} else if err != nil {
		log.Fatalln(err)
	} else {
		return hres, nil
	}

	// Create a expiration time variable so that SETEX wouldn't cause millisecond delay from individual SETs
	expireAt := time.Duration(exp) * time.Minute
	m := map[string]interface{}{
		"short":                 h,
		"full":                  url,
		"created_at":            time.Now().String(),
		"expiration_in_minutes": expireAt.String(),
	}

	batch := func(tx *redis.Tx) error {
		// Queue HMSETs
		err = tx.HMSet(h, m).Err()
		if err != nil {
			return err
		}
		err = tx.Expire(h, expireAt).Err()
		if err != nil {
			return err
		}

		err = tx.HMSet(url, m).Err()
		if err != nil {
			return err
		}
		err = tx.Expire(url, expireAt).Err()
		if err != nil {
			return err
		}

		// Exec
		//_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
		//	// Dunno what to do here.
		//	return nil
		//})

		return nil
	}

	err = r.Cli.Watch(batch, h, url)
	//log.Println("Watch", err)
	if err == redis.TxFailedErr {
		return "", err
	}

	//log.Println("watched", err)
	return h, nil
}

// ShortLinkInfo gets ShortlinkDetailKey from redis
func (r *RedisCli) ShortLinkInfo(eid string) (*DetailInfo, error) {

	hres, err := r.Cli.HGetAll(eid).Result()
	//r.Cli.HMGet(eid, "short", "full", "created_at", "expiration_in_minutes").Result() //r.Cli.HGetAll(eid).Result()
	//log.Println("hmget", hres, err)
	// Check if HMGET result is nil
	if err == redis.Nil {
		return nil, myerror.StatusError{
			Code: 404,
			Err:  errors.New("unknown short URL"),
		}
	} else if err != nil {
		return nil, err
	} else if len(hres) == 0 {
		return nil, myerror.StatusError{
			Code: 404,
			Err:  errors.New("unknown short URL"),
		}
	}

	exp, err := time.ParseDuration(hres["expiration_in_minutes"])
	if err != nil {
		log.Fatalln("ParseDuration", err)
	}

	sli := &DetailInfo{
		Short:               hres["short"],
		Full:                hres["full"],
		CreatedAt:           hres["created_at"],
		ExpirationInMinutes: exp / time.Minute,
	}

	return sli, nil
}
