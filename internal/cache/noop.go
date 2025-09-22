package cache

import (
	"context"
	"time"
)

// NoOpCache is a no-operation cache that does nothing (for when caching is disabled)
type NoOpCache struct{}

// NewNoOpCache creates a new no-op cache
func NewNoOpCache() *NoOpCache {
	return &NoOpCache{}
}

// Get always returns cache miss
func (n *NoOpCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrCacheMiss
}

// Set does nothing
func (n *NoOpCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

// Delete does nothing
func (n *NoOpCache) Delete(ctx context.Context, key string) error {
	return nil
}

// DeletePattern does nothing
func (n *NoOpCache) DeletePattern(ctx context.Context, pattern string) error {
	return nil
}

// Exists always returns false
func (n *NoOpCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// SetNX always returns false
func (n *NoOpCache) SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error) {
	return false, nil
}

// Close does nothing
func (n *NoOpCache) Close() error {
	return nil
}