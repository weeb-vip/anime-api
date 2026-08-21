package anime_staff

import (
	"context"
	"time"

	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
)

type AnimeStaffRepositoryImpl interface {
	FindStaffByID(ctx context.Context, id string) (*AnimeStaff, error)
}

type AnimeStaffRepository struct {
	db *db.DB
}

func NewAnimeStaffRepository(db *db.DB) AnimeStaffRepositoryImpl {
	return &AnimeStaffRepository{db: db}
}

func (a *AnimeStaffRepository) FindStaffByID(ctx context.Context, id string) (*AnimeStaff, error) {
	startTime := time.Now()

	var staff AnimeStaff
	err := a.db.DB.WithContext(ctx).Where("id = ?", id).First(&staff).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_staff",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_staff",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return &staff, nil
}
