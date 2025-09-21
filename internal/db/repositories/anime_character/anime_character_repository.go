package anime_character

import (
	"context"
	"time"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
)

type AnimeCharacterRepositoryImpl interface {
	FindAnimearacterById(ctx context.Context, id string) (*AnimeCharacter, error)
}

type AnimeCharacterRepository struct {
	db *db.DB
}

func NewAnimeCharacterRepository(db *db.DB) AnimeCharacterRepositoryImpl {
	return &AnimeCharacterRepository{db: db}
}

func (a *AnimeCharacterRepository) FindAnimearacterById(ctx context.Context, id string) (*AnimeCharacter, error) {
	startTime := time.Now()

	var animeCharacter AnimeCharacter
	err := a.db.DB.WithContext(ctx).Where("id = ?", id).First(&animeCharacter).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_characters",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_characters",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return &animeCharacter, nil
}
