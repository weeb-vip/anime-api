package cache

import (
	"context"
	"strings"
)

// CacheCoordinator handles coordinated cache invalidation between related entities
type CacheCoordinator struct {
	cache *CacheService
}

// NewCacheCoordinator creates a new cache coordinator
func NewCacheCoordinator(cache *CacheService) *CacheCoordinator {
	return &CacheCoordinator{
		cache: cache,
	}
}

// InvalidateAnimeAndRelated invalidates anime cache and related episode cache
func (c *CacheCoordinator) InvalidateAnimeAndRelated(ctx context.Context, animeID string) error {
	// Invalidate specific anime
	animeKey := c.cache.GetKeyBuilder().AnimeByID(animeID)
	_ = c.cache.Delete(ctx, animeKey)

	// Invalidate episodes for this anime
	episodeKey := c.cache.GetKeyBuilder().EpisodesByAnimeID(animeID)
	_ = c.cache.Delete(ctx, episodeKey)

	// Invalidate all season caches (since they might contain this anime)
	seasonPattern := c.cache.GetKeyBuilder().AnimeBySeasonPattern("*")
	seasonPattern = strings.Replace(seasonPattern, ":*", "*", 1)
	_ = c.cache.DeletePattern(ctx, seasonPattern)

	// Invalidate lists (top rated, popular, newest)
	listPattern := c.cache.GetKeyBuilder().AnimePattern()
	_ = c.cache.DeletePattern(ctx, listPattern)

	return nil
}

// InvalidateEpisodesOnly invalidates only episode cache for a specific anime
func (c *CacheCoordinator) InvalidateEpisodesOnly(ctx context.Context, animeID string) error {
	// Invalidate episodes for this anime
	episodeKey := c.cache.GetKeyBuilder().EpisodesByAnimeID(animeID)
	_ = c.cache.Delete(ctx, episodeKey)

	return nil
}

// InvalidateSeasonCaches invalidates all season-related caches
func (c *CacheCoordinator) InvalidateSeasonCaches(ctx context.Context) error {
	// Invalidate all season caches
	seasonPattern := c.cache.GetKeyBuilder().AnimeBySeasonPattern("*")
	seasonPattern = strings.Replace(seasonPattern, ":*", "*", 1)
	_ = c.cache.DeletePattern(ctx, seasonPattern)

	return nil
}

// InvalidateListCaches invalidates ranking list caches (top rated, popular, newest)
func (c *CacheCoordinator) InvalidateListCaches(ctx context.Context) error {
	// Get base pattern and create specific patterns for lists
	basePattern := c.cache.GetKeyBuilder().AnimePattern()
	baseKey := basePattern[:len(basePattern)-1] // Remove the trailing "*"

	// Invalidate specific list patterns
	patterns := []string{
		baseKey + "top_rated:*",
		baseKey + "most_popular:*",
		baseKey + "newest:*",
	}

	for _, pattern := range patterns {
		_ = c.cache.DeletePattern(ctx, pattern)
	}

	return nil
}

// InvalidateAllAnime invalidates all anime-related caches
func (c *CacheCoordinator) InvalidateAllAnime(ctx context.Context) error {
	// Invalidate all anime caches
	animePattern := c.cache.GetKeyBuilder().AnimePattern()
	_ = c.cache.DeletePattern(ctx, animePattern)

	// Invalidate all season caches
	seasonPattern := c.cache.GetKeyBuilder().AnimeBySeasonPattern("*")
	seasonPattern = strings.Replace(seasonPattern, ":*", "*", 1)
	_ = c.cache.DeletePattern(ctx, seasonPattern)

	return nil
}

// InvalidateAllEpisodes invalidates all episode-related caches
func (c *CacheCoordinator) InvalidateAllEpisodes(ctx context.Context) error {
	// Invalidate all episode caches
	episodePattern := c.cache.GetKeyBuilder().EpisodePattern()
	episodePattern = strings.Replace(episodePattern, ":*", "*", 1)
	_ = c.cache.DeletePattern(ctx, episodePattern)

	return nil
}

// ClearAllCache clears all cache entries (use with caution)
func (c *CacheCoordinator) ClearAllCache(ctx context.Context) error {
	// This is a nuclear option - clears everything
	// Be very careful using this in production
	return c.cache.DeletePattern(ctx, "*")
}