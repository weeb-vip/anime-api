package anime

import (
	"context"
	"time"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
)

type RECORD_TYPE string

type AnimeEpisodeRepositoryImpl interface {
	Upsert(ctx context.Context, anime *AnimeEpisode) error
	Delete(ctx context.Context, anime *AnimeEpisode) error
	FindByAnimeID(ctx context.Context, animeID string) ([]*AnimeEpisode, error)
}

type AnimeEpisodeRepository struct {
	db *db.DB
}

func NewAnimeEpisodeRepository(db *db.DB) AnimeEpisodeRepositoryImpl {
	return &AnimeEpisodeRepository{db: db}
}

func (a *AnimeEpisodeRepository) Upsert(ctx context.Context, episode *AnimeEpisode) error {
	startTime := time.Now()

	err := a.db.DB.WithContext(ctx).Save(episode).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_episodes",
			Method:  metrics_lib.DatabaseMetricMethodInsert,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_episodes",
		Method:  metrics_lib.DatabaseMetricMethodInsert,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return nil
}

func (a *AnimeEpisodeRepository) Delete(ctx context.Context, episode *AnimeEpisode) error {
	startTime := time.Now()

	err := a.db.DB.WithContext(ctx).Delete(episode).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_episodes",
			Method:  metrics_lib.DatabaseMetricMethodDelete,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_episodes",
		Method:  metrics_lib.DatabaseMetricMethodDelete,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return nil
}

func (a *AnimeEpisodeRepository) FindByAnimeID(ctx context.Context, animeID string) ([]*AnimeEpisode, error) {
	startTime := time.Now()

	var episodes []*AnimeEpisode
	err := a.db.DB.WithContext(ctx).Where("anime_id = ?", animeID).Order("episode ASC").Find(&episodes).Error
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

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_episodes",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return episodes, nil
}
