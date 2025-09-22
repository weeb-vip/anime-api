package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/weeb-vip/anime-api/metrics"
)

// CacheService provides high-level caching operations with JSON serialization
type CacheService struct {
	cache      Cache
	keyBuilder *CacheKeyBuilder
}

// NewCacheService creates a new cache service
func NewCacheService(cache Cache) *CacheService {
	return &CacheService{
		cache:      cache,
		keyBuilder: GetKeyBuilder(),
	}
}

// GetJSON retrieves and unmarshals JSON data from cache
func (c *CacheService) GetJSON(ctx context.Context, key string, dest interface{}) error {
	startTime := time.Now()

	data, err := c.cache.Get(ctx, key)
	if err != nil {
		if err == ErrCacheMiss {
			metrics.GetAppMetrics().DatabaseMetric(
				float64(time.Since(startTime).Milliseconds()),
				"cache",
				"get",
				"miss",
			)
		} else {
			metrics.GetAppMetrics().DatabaseMetric(
				float64(time.Since(startTime).Milliseconds()),
				"cache",
				"get",
				metrics.Error,
			)
		}
		return err
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"cache",
			"get",
			metrics.Error,
		)
		return fmt.Errorf("cache unmarshal error: %w", err)
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"cache",
		"get",
		"hit",
	)
	return nil
}

// SetJSON marshals and stores JSON data in cache
func (c *CacheService) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	startTime := time.Now()

	data, err := json.Marshal(value)
	if err != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"cache",
			"set",
			metrics.Error,
		)
		return fmt.Errorf("cache marshal error: %w", err)
	}

	err = c.cache.Set(ctx, key, data, ttl)
	if err != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"cache",
			"set",
			metrics.Error,
		)
		return err
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"cache",
		"set",
		metrics.Success,
	)
	return nil
}

// Delete removes a key from cache
func (c *CacheService) Delete(ctx context.Context, key string) error {
	startTime := time.Now()

	err := c.cache.Delete(ctx, key)

	result := metrics.Success
	if err != nil {
		result = metrics.Error
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"cache",
		"delete",
		result,
	)
	return err
}

// DeletePattern removes all keys matching a pattern
func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	startTime := time.Now()

	err := c.cache.DeletePattern(ctx, pattern)

	result := metrics.Success
	if err != nil {
		result = metrics.Error
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"cache",
		"delete_pattern",
		result,
	)
	return err
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	return c.cache.Exists(ctx, key)
}

// GetKeyBuilder returns the cache key builder
func (c *CacheService) GetKeyBuilder() *CacheKeyBuilder {
	return c.keyBuilder
}

// Close closes the cache connection
func (c *CacheService) Close() error {
	return c.cache.Close()
}