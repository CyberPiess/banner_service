package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Config struct {
	Addres        string
	RedisPassword string
}

func NewRedis(cfg Config) (*redis.Client, error) {

	var client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addres,
		Password: cfg.RedisPassword,
	})

	if client == nil {
		return nil, fmt.Errorf("redis is not running")
	}

	return client, nil
}
