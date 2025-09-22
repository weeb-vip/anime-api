package anime

import (
	"context"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
	"gorm.io/gorm"
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
	err := a.db.DB.WithContext(ctx).
		Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
			return db.Order("episode ASC")
		}).
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

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}