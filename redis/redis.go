package redis

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(url string) *redis.Client {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	opt.DB = 4
	client := redis.NewClient(opt)
	return client
}
