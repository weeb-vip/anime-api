package anime_schedule

import (
	"github.com/weeb-vip/anime-api/internal/db"
)

type AnimeScheduleRepositoryImpl interface {
	FindByAnimeID(animeID string) (*AnimeSchedule, error)
}

type AnimeScheduleRepository struct {
	db *db.DB
}

func NewAnimeScheduleRepository(db *db.DB) AnimeScheduleRepositoryImpl {
	return &AnimeScheduleRepository{db: db}
}

func (r *AnimeScheduleRepository) FindByAnimeID(animeID string) (*AnimeSchedule, error) {
	var schedule AnimeSchedule
	err := r.db.DB.Where("anime_id = ?", animeID).First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}
