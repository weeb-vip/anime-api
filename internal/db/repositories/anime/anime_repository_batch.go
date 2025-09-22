package anime

import (
	"context"
	anime_episode "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
	"time"
)

// FindByIDs fetches multiple anime by their IDs in a single query
func (a *AnimeRepository) FindByIDs(ctx context.Context, ids []string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

// FindByIDsWithEpisodes fetches multiple anime by their IDs with episodes preloaded in a single query
func (a *AnimeRepository) FindByIDsWithEpisodes(ctx context.Context, ids []string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime

	// First, fetch all anime
	err := a.db.DB.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&animes).Error

	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Then, batch fetch all episodes for these anime in one query
	var episodes []anime_episode.AnimeEpisode
	err = a.db.DB.WithContext(ctx).
		Where("anime_id IN ?", ids).
		Order("anime_id ASC, episode ASC").
		Find(&episodes).Error

	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_episodes",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Group episodes by anime_id
	episodeMap := make(map[string][]*anime_episode.AnimeEpisode)
	for _, episode := range episodes {
		if episode.AnimeID != nil {
			ep := episode // Create a copy to avoid pointer issues
			episodeMap[*episode.AnimeID] = append(episodeMap[*episode.AnimeID], &ep)
		}
	}

	// Assign episodes to their respective anime
	for i, anime := range animes {
		if eps, exists := episodeMap[anime.ID]; exists {
			animes[i].AnimeEpisodes = eps
		}
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}