package anime_news

import (
	"github.com/weeb-vip/anime-api/internal/db"
)

type AnimeNewsRepositoryImpl interface {
	FindByAnimeID(animeID string) ([]AnimeNews, error)
}

type AnimeNewsRepository struct {
	db *db.DB
}

func NewAnimeNewsRepository(db *db.DB) AnimeNewsRepositoryImpl {
	return &AnimeNewsRepository{db: db}
}

func (r *AnimeNewsRepository) FindByAnimeID(animeID string) ([]AnimeNews, error) {
	var news []AnimeNews
	err := r.db.DB.Where("anime_id = ?", animeID).Order("published_date DESC, created_at DESC").Find(&news).Error
	return news, err
}
