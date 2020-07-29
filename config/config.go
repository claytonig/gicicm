package config

import (
	"fmt"
	"os"

	"gicicm/logger"
)

// DbConfig contains the database configuration details.
type DbConfig struct {
	Host   string // DB_HOST
	Port   string // DB_PORT
	User   string // DB_USER
	Pass   string // DB_PASS
	DBName string // DB_NAME
	DBType string // DB_TYPE
}

// CacheConfig contains the cache configuration details.
type CacheConfig struct {
	Host string // CACHE_HOST
}

// Config contains configuration details for gicicm to start
type Config struct {
	Database   DbConfig
	Cache      CacheConfig
	SigningKey string
}

// GetConfig returns an instance of config
func GetConfig() *Config {

	dbConf := DbConfig{
		Host:   mustGetEnv("DB_HOST"),
		Port:   mustGetEnv("DB_PORT"),
		User:   mustGetEnv("DB_USER"),
		Pass:   mustGetEnv("DB_PASS"),
		DBName: mustGetEnv("DB_NAME"),
		DBType: mustGetEnv("DB_TYPE"),
	}

	cacheConf := CacheConfig{
		Host: mustGetEnv("CACHE_HOST"),
	}

	return &Config{
		Database:   dbConf,
		Cache:      cacheConf,
		SigningKey: mustGetEnv("SIGNING_KEY"),
	}
}

// mustGetEnv returns the value of an env variable
// panics if not set.
func mustGetEnv(env string) string {
	value := os.Getenv(env)
	if value == "" {
		logger.Log().Panic(fmt.Sprintf("The environment variable %s is missing.", env))
	}
	return value
}
