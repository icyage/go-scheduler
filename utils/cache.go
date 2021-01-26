package utils

import (
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

func GetCache() *redis.Client {

	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.GetString("redis_host"),
			Password: cfg.GetString("redis_password"), // no password set
			DB:       cfg.GetInt("redis_db"),          // use default DB
		})
		return client
	} else {
		return client
	}
}
