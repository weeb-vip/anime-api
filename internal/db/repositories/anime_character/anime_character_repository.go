package anime_character

import (
	"context"
	"time"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/metrics"
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
		metrics.GetAppMetrics().DatabaseSince(startTime, "anime_characters", metrics.MethodSelect, metrics.Error)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseSince(startTime, "anime_characters", metrics.MethodSelect, metrics.Success)
	return &animeCharacter, nil
}
