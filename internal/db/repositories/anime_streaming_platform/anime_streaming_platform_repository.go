package anime_streaming_platform

import (
	"github.com/weeb-vip/anime-api/internal/db"
)

type AnimeStreamingPlatformRepositoryImpl interface {
	FindByAnimeID(animeID string) ([]AnimeStreamingPlatform, error)
}

type AnimeStreamingPlatformRepository struct {
	db *db.DB
}

func NewAnimeStreamingPlatformRepository(db *db.DB) AnimeStreamingPlatformRepositoryImpl {
	return &AnimeStreamingPlatformRepository{db: db}
}

func (r *AnimeStreamingPlatformRepository) FindByAnimeID(animeID string) ([]AnimeStreamingPlatform, error) {
	var platforms []AnimeStreamingPlatform
	err := r.db.DB.Where("anime_id = ?", animeID).Find(&platforms).Error
	return platforms, err
}
