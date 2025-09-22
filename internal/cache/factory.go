package cache

import (
	"fmt"

	"github.com/weeb-vip/anime-api/config"
)

// Factory creates cache instances based on configuration
func NewCache(cfg config.Config) (Cache, error) {
	if !cfg.RedisConfig.Enabled {
		return NewNoOpCache(), nil
	}

	redisCache, err := NewRedisCache(cfg.RedisConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis cache: %w", err)
	}

	return redisCache, nil
}

// GetKeyBuilder returns a cache key builder
func GetKeyBuilder() *CacheKeyBuilder {
	return NewCacheKeyBuilder("anime-api")
}