package config

import (
	"log"
	"os"
	"strconv"
)

// RedisEnv contains redis addr, password and db
type RedisEnv struct {
	Addr     string
	Password string
	DB       int
}

// AppEnv contains app addr
type AppEnv struct {
	Addr string
}

// GetRedisEnv returns a RedisEnv
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

// GetAppEnv returns a AppEnv
func GetAppEnv() AppEnv {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = "0.0.0.0:8080"
	}
	return AppEnv{Addr: addr}
}
