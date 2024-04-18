package main

import (
	"os"

	"github.com/go-redis/redis"
)

// NOTE: user redis to store this DB
var DB map[string]*Bucket

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func main() {
}
