package redis_cache

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type RedisCache struct {
	client *redis.Client
	logger logger
}

func NewBannerCache(client *redis.Client, logger logger) *RedisCache {
	return &RedisCache{client: client, logger: logger}
}

func (rc *RedisCache) IfCacheExists(key string) (int64, error) {

	result, err := rc.client.Exists(key).Result()
	if err != nil {
		rc.logger.WithFields(logrus.Fields{
			"package":  "redis_cache",
			"function": "IfCacheExists",
			"error":    err,
		}).Error("Error searchng key in Redis")
	}

	return result, err
}

func (rc *RedisCache) AddToCache(key string, redisDTO RedisEntity) error {

	_, err := rc.client.Set(key, redisDTO.Content, time.Minute*5).Result()
	if err != nil {
		rc.logger.WithFields(logrus.Fields{
			"package":  "redis_cache",
			"function": "AddToCache",
			"error":    err,
		}).Error("Error adding key to Redis")
	}

	return err
}

func (rc *RedisCache) GetFromCache(key string) (RedisEntity, error) {

	var redisDTO RedisEntity

	err := rc.client.Get(key).Scan(&redisDTO.Content)
	if err != nil {
		rc.logger.WithFields(logrus.Fields{
			"package":  "redis_cache",
			"function": "GetFromCache",
			"error":    err,
		}).Error("Error getting key from Redis")
	}

	return redisDTO, err
}
