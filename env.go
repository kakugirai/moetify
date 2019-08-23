package main

import (
	"log"
	"os"
	"strconv"
)

type RedisEnv struct {
	RS RedisStorage
}

type AppEnv struct {
	Addr string
}

func getRedisEnv() *RedisEnv {
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	passwd := os.Getenv("APP_REDIS_PASSWD")
	if passwd == "" {
		passwd = ""
	}

	dbS := os.Getenv("APP_REDIS_DB")
	if dbS == "" {
		dbS = "0"
	}
	db, err := strconv.Atoi(dbS)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connect to redis (addr: %s, password: %s, db: %d)", addr, passwd, db)

	r := NewRedisCli(addr, passwd, db)
	return &RedisEnv{RS: r}
}

func getAppEnv() AppEnv {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = "0.0.0.0:8080"
	}
	return AppEnv{Addr: addr}
}
