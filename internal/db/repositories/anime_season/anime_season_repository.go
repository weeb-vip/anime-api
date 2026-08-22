package anime_season

import (
	"context"
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
		metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodSelect, metrics.Success)
	return animeSeasons, nil
}

func (r *AnimeSeasonRepository) FindBySeason(ctx context.Context, season string) ([]*AnimeSeason, error) {
	startTime := time.Now()

	var animeSeasons []*AnimeSeason
	err := r.db.DB.Where("season = ?", season).Find(&animeSeasons).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodSelect, metrics.Success)
	return animeSeasons, nil
}

func (r *AnimeSeasonRepository) Create(ctx context.Context, animeSeason *AnimeSeason) error {
	startTime := time.Now()

	err := r.db.DB.Create(animeSeason).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodInsert, metrics.Error)
		return err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodInsert, metrics.Success)
	return nil
}

func (r *AnimeSeasonRepository) Update(ctx context.Context, animeSeason *AnimeSeason) error {
	startTime := time.Now()

	err := r.db.DB.Save(animeSeason).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodUpdate, metrics.Error)
		return err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodUpdate, metrics.Success)
	return nil
}

func (r *AnimeSeasonRepository) Delete(ctx context.Context, id string) error {
	startTime := time.Now()

	err := r.db.DB.Delete(&AnimeSeason{}, "id = ?", id).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodDelete, metrics.Error)
		return err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "anime_seasons", metrics.MethodDelete, metrics.Success)
	return nil
}