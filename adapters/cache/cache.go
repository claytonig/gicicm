package cache

import (
	"github.com/go-redis/redis"
	"time"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, duration time.Duration) (string, error)
	Del(key string) error
}

type cache struct {
	cacheConn redis.Client
}
