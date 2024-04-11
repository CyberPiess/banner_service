package redis_cache

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisCache struct {
	client *redis.Client
}

func NewBannerCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (rc *RedisCache) AddToCache(key string, redisDTO RedisEntity) error {

	_, err := rc.client.Set(key, redisDTO.Content, time.Minute*5).Result()

	return err
}

func (rc *RedisCache) GetFromCache(key string) (RedisEntity, error) {

	var redisDTO RedisEntity

	err := rc.client.Get(key).Scan(&redisDTO.Content)

	return redisDTO, err
}

func (rc *RedisCache) DeleteFromCache(key string) error {
	err := rc.client.Del(key).Err()

	return err
}
