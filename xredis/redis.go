package xredis

import "github.com/go-redis/redis"

func New(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Addr,
	})
}
