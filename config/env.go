package config

import (
	"log"
	"os"
	"strconv"
)

type RedisEnv struct {
	Addr     string
	Password string
	DB       int
}

type AppEnv struct {
	Addr string
}

func GetRedisEnv() RedisEnv {
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

	return RedisEnv{
		addr,
		passwd,
		db,
	}
}

func GetAppEnv() AppEnv {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = "0.0.0.0:8080"
	}
	return AppEnv{Addr: addr}
}
