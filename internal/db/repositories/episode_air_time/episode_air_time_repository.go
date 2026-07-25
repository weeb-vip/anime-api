package episode_air_time

import (
	"github.com/weeb-vip/anime-api/internal/db"
)

type EpisodeAirTimeRepositoryImpl interface {
	FindByAnimeID(animeID string) ([]EpisodeAirTime, error)
	FindByAnimeIDAndEpisode(animeID string, episodeNumber int) ([]EpisodeAirTime, error)
	FindSubTimeByAnimeIDAndEpisode(animeID string, episodeNumber int) (*EpisodeAirTime, error)
}

type EpisodeAirTimeRepository struct {
	db *db.DB
}

func NewEpisodeAirTimeRepository(db *db.DB) EpisodeAirTimeRepositoryImpl {
	return &EpisodeAirTimeRepository{db: db}
}

func (r *EpisodeAirTimeRepository) FindByAnimeID(animeID string) ([]EpisodeAirTime, error) {
	var airTimes []EpisodeAirTime
	err := r.db.DB.Where("anime_id = ?", animeID).
		Order("episode_number ASC, air_type ASC").
		Find(&airTimes).Error
	return airTimes, err
}

func (r *EpisodeAirTimeRepository) FindByAnimeIDAndEpisode(animeID string, episodeNumber int) ([]EpisodeAirTime, error) {
	var airTimes []EpisodeAirTime
	err := r.db.DB.Where("anime_id = ? AND episode_number = ?", animeID, episodeNumber).
		Order("air_type ASC").
		Find(&airTimes).Error
	return airTimes, err
}

func (r *EpisodeAirTimeRepository) FindSubTimeByAnimeIDAndEpisode(animeID string, episodeNumber int) (*EpisodeAirTime, error) {
	var airTime EpisodeAirTime
	err := r.db.DB.Where("anime_id = ? AND episode_number = ? AND air_type = 'sub'", animeID, episodeNumber).
		First(&airTime).Error
	if err != nil {
		return nil, err
	}
	return &airTime, nil
}
