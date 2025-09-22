package cache

import "errors"

var (
	// ErrCacheMiss is returned when a cache key is not found
	ErrCacheMiss = errors.New("cache miss")

	// ErrCacheDisabled is returned when caching is disabled
	ErrCacheDisabled = errors.New("cache disabled")
)