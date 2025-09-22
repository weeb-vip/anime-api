package cache

import (
	"context"
	"time"

	"github.com/weeb-vip/anime-api/config"
)

// Cache defines the interface for caching operations
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in cache with TTL
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// DeletePattern removes all keys matching a pattern
	DeletePattern(ctx context.Context, pattern string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// SetNX sets a value only if the key doesn't exist (for locking)
	SetNX(ctx context.Context, key string, value []byte, ttl time.Duration) (bool, error)

	// Close closes the cache connection
	Close() error
}

// CacheKeyBuilder helps build consistent cache keys
type CacheKeyBuilder struct {
	prefix string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{prefix: prefix}
}

// AnimeByID builds cache key for anime by ID
func (c *CacheKeyBuilder) AnimeByID(id string) string {
	return c.prefix + ":anime:id:" + id
}

// AnimeBySeasonPattern builds cache key pattern for anime by season
func (c *CacheKeyBuilder) AnimeBySeasonPattern(season string) string {
	return c.prefix + ":anime:season:" + season + ":*"
}

// AnimeBySeasonWithFields builds cache key for anime by season with specific fields
func (c *CacheKeyBuilder) AnimeBySeasonWithFields(season string, fields []string) string {
	if len(fields) == 0 {
		return c.prefix + ":anime:season:" + season + ":all"
	}

	key := c.prefix + ":anime:season:" + season + ":fields:"
	for i, field := range fields {
		if i > 0 {
			key += ","
		}
		key += field
	}
	return key
}

// EpisodesByAnimeID builds cache key for episodes by anime ID
func (c *CacheKeyBuilder) EpisodesByAnimeID(animeID string) string {
	return c.prefix + ":episodes:anime:" + animeID
}

// EpisodeByID builds cache key for episode by ID
func (c *CacheKeyBuilder) EpisodeByID(id string) string {
	return c.prefix + ":episode:id:" + id
}

// AnimePattern builds pattern for all anime cache keys
func (c *CacheKeyBuilder) AnimePattern() string {
	return c.prefix + ":anime:*"
}

// EpisodePattern builds pattern for all episode cache keys
func (c *CacheKeyBuilder) EpisodePattern() string {
	return c.prefix + ":episode*"
}

// AnimeByIDPattern builds pattern for anime invalidation by ID
func (c *CacheKeyBuilder) AnimeByIDPattern(animeID string) string {
	return c.prefix + ":*anime*:" + animeID + "*"
}

// TTL helper functions that use configuration values
func GetAnimeDataTTL(cfg config.RedisConfig) time.Duration {
	return time.Duration(cfg.AnimeDataTTLMinutes) * time.Minute
}

func GetEpisodeTTL(cfg config.RedisConfig) time.Duration {
	return time.Duration(cfg.EpisodeTTLMinutes) * time.Minute
}

func GetSeasonTTL(cfg config.RedisConfig) time.Duration {
	return time.Duration(cfg.SeasonTTLMinutes) * time.Minute
}

func GetLockTTL(cfg config.RedisConfig) time.Duration {
	return time.Duration(cfg.LockTTLSeconds) * time.Second
}