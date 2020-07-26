package cache

import (
	"time"

	"gicicm/config"
	"gicicm/logger"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, duration time.Duration) (string, error)
	Del(key string) error
}

type cache struct {
	cacheConn *redis.Client
}

// NewCache returns an instance of newCache
func NewCache(config *config.Config) Cache {
	cacheConn := newCacheConnection(config)
	return &cache{
		cacheConn: cacheConn,
	}
}

// newCacheConnection - Initializes cache connection
func newCacheConnection(config *config.Config) *redis.Client {
	cacheConn := redis.NewClient(&redis.Options{
		Addr:        config.Cache.Host,
		Password:    "",
		DB:          0,
		ReadTimeout: time.Second,
	})
	if cacheConn == nil {
		logger.Log().Fatal("unable to connect to redis", zap.String("host", config.Cache.Host))
	}
	return cacheConn
}

// Get - Get value from redis
func (c *cache) Get(key string) (string, error) {
	data, err := c.cacheConn.Get(key).Result()
	if err != nil {
		logger.Log().Error("Error while fetching data from redis", zap.String("key", key), zap.Error(err))
	}
	return data, err
}

// Set - Set value to redis
func (c *cache) Set(key string, value string, duration time.Duration) (string, error) {
	result, err := c.cacheConn.Set(key, value, duration).Result()
	if err != nil {
		logger.Log().Error("Error while storing data to redis", zap.String("key", key), zap.Error(err))
	}
	return result, err
}

// Del - Delete the key
func (c *cache) Del(key string) error {
	err := c.cacheConn.Del(key).Err()
	if err != nil {
		logger.Log().Error("Error while deleting key", zap.String("key", key), zap.Error(err))
		return err
	}
	return nil
}
