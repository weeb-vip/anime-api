package anime

import (
	"context"
	"time"
	"github.com/weeb-vip/anime-api/internal/cache"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
)

type RECORD_TYPE string

type AnimeEpisodeRepositoryImpl interface {
	Upsert(ctx context.Context, anime *AnimeEpisode) error
	Delete(ctx context.Context, anime *AnimeEpisode) error
	FindByAnimeID(ctx context.Context, animeID string) ([]*AnimeEpisode, error)
}

type AnimeEpisodeRepository struct {
	db    *db.DB
	cache *cache.CacheService
}

func NewAnimeEpisodeRepository(db *db.DB) AnimeEpisodeRepositoryImpl {
	return &AnimeEpisodeRepository{db: db, cache: nil}
}

func NewAnimeEpisodeRepositoryWithCache(db *db.DB, cacheService *cache.CacheService) AnimeEpisodeRepositoryImpl {
	return &AnimeEpisodeRepository{db: db, cache: cacheService}
}

func (a *AnimeEpisodeRepository) Upsert(ctx context.Context, episode *AnimeEpisode) error {
	startTime := time.Now()

	err := a.db.DB.WithContext(ctx).Save(episode).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"anime_episodes",
			"insert",
			metrics.Error,
		)
		return err
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"anime_episodes",
		"insert",
		metrics.Success,
	)

	// Invalidate cache if available
	if a.cache != nil && episode.AnimeID != nil {
		coordinator := cache.NewCacheCoordinator(a.cache)
		_ = coordinator.InvalidateEpisodesOnly(ctx, *episode.AnimeID)
	}

	return nil
}

func (a *AnimeEpisodeRepository) Delete(ctx context.Context, episode *AnimeEpisode) error {
	startTime := time.Now()

	err := a.db.DB.WithContext(ctx).Delete(episode).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"anime_episodes",
			"delete",
			metrics.Error,
		)
		return err
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"anime_episodes",
		"delete",
		metrics.Success,
	)

	// Invalidate cache if available
	if a.cache != nil && episode.AnimeID != nil {
		coordinator := cache.NewCacheCoordinator(a.cache)
		_ = coordinator.InvalidateEpisodesOnly(ctx, *episode.AnimeID)
	}

	return nil
}

func (a *AnimeEpisodeRepository) FindByAnimeID(ctx context.Context, animeID string) ([]*AnimeEpisode, error) {
	// Try cache first if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().EpisodesByAnimeID(animeID)
		var episodes []*AnimeEpisode
		err := a.cache.GetJSON(ctx, key, &episodes)
		if err == nil {
			return episodes, nil
		}
		// Continue to database if cache miss or error
	}

	startTime := time.Now()
	var episodes []*AnimeEpisode
	err := a.db.DB.WithContext(ctx).Where("anime_id = ?", animeID).Order("episode ASC").Find(&episodes).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"anime_episodes",
			"select",
			metrics.Error,
		)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"anime_episodes",
		"select",
		metrics.Success,
	)

	// Store in cache if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().EpisodesByAnimeID(animeID)
		_ = a.cache.SetJSON(ctx, key, episodes, cache.EpisodeTTL)
	}

	return episodes, nil
}
