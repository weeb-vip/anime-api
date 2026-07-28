package anime_fanart

import (
	"github.com/weeb-vip/anime-api/internal/db"
)

type AnimeFanartRepositoryImpl interface {
	FindByAnimeID(animeID string) ([]Fanart, error)
}

type AnimeFanartRepository struct {
	db *db.DB
}

func NewAnimeFanartRepository(db *db.DB) AnimeFanartRepositoryImpl {
	return &AnimeFanartRepository{db: db}
}

func (r *AnimeFanartRepository) FindByAnimeID(animeID string) ([]Fanart, error) {
	var fanart []Fanart
	err := r.db.DB.Where("anime_id = ?", animeID).Order("created_at DESC").Find(&fanart).Error
	return fanart, err
}
