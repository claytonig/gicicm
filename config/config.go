package conf

// dbConfig contains the database configuration details.
type dbConfig struct {
}

// cacheConfig contains the cache configuration details.
type cacheConfig struct {
}

// Config contains configuration details for gicicm to start
type Config struct {
	Database dbConfig
	Cache    cacheConfig
}

// GetConfig returns an instance of config
// if there is an error initializing any of
// the config variables an error is returned.
func GetConfig() (*Config, error) {
	return &Config{}, nil
}
