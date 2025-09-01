package anime_season

import (
	"context"
	metrics_lib "github.com/TempMee/go-metrics-lib"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
	"time"
)

type AnimeSeasonRepositoryImpl interface {
	FindByAnimeID(ctx context.Context, animeID string) ([]*AnimeSeason, error)
	FindBySeason(ctx context.Context, season string) ([]*AnimeSeason, error)
	Create(ctx context.Context, animeSeason *AnimeSeason) error
	Update(ctx context.Context, animeSeason *AnimeSeason) error
	Delete(ctx context.Context, id string) error
}

type AnimeSeasonRepository struct {
	db *db.DB
}

func NewAnimeSeasonRepository(db *db.DB) AnimeSeasonRepositoryImpl {
	return &AnimeSeasonRepository{db: db}
}

func (r *AnimeSeasonRepository) FindByAnimeID(ctx context.Context, animeID string) ([]*AnimeSeason, error) {
	startTime := time.Now()

	var animeSeasons []*AnimeSeason
	err := r.db.DB.Where("anime_id = ?", animeID).Find(&animeSeasons).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animeSeasons, nil
}

func (r *AnimeSeasonRepository) FindBySeason(ctx context.Context, season string) ([]*AnimeSeason, error) {
	startTime := time.Now()

	var animeSeasons []*AnimeSeason
	err := r.db.DB.Where("season = ?", season).Find(&animeSeasons).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animeSeasons, nil
}

func (r *AnimeSeasonRepository) Create(ctx context.Context, animeSeason *AnimeSeason) error {
	startTime := time.Now()

	err := r.db.DB.Create(animeSeason).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodInsert,
			Result:  metrics_lib.Error,
		})
		return err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodInsert,
		Result:  metrics_lib.Success,
	})
	return nil
}

func (r *AnimeSeasonRepository) Update(ctx context.Context, animeSeason *AnimeSeason) error {
	startTime := time.Now()

	err := r.db.DB.Save(animeSeason).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodUpdate,
			Result:  metrics_lib.Error,
		})
		return err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodUpdate,
		Result:  metrics_lib.Success,
	})
	return nil
}

func (r *AnimeSeasonRepository) Delete(ctx context.Context, id string) error {
	startTime := time.Now()

	err := r.db.DB.Delete(&AnimeSeason{}, "id = ?", id).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodDelete,
			Result:  metrics_lib.Error,
		})
		return err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodDelete,
		Result:  metrics_lib.Success,
	})
	return nil
}