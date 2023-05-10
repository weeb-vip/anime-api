package anime

import "github.com/weeb-vip/anime-api/internal/db"

type RECORD_TYPE string

type AnimeRepositoryImpl interface {
	FindAll() ([]*Anime, error)
	FindById(id string) (*Anime, error)
	FindByName(name string) ([]*Anime, error)
	FindByType(recordType RECORD_TYPE) ([]*Anime, error)
	FindByStatus(status string) ([]*Anime, error)
	FindBySource(source string) ([]*Anime, error)
	FindByGenre(genre string) ([]*Anime, error)
	FindByStudio(studio string) ([]*Anime, error)
	FindByLicensors(licensors string) ([]*Anime, error)
	FindByRating(rating string) ([]*Anime, error)
	FindByYear(year int) ([]*Anime, error)
	FindBySeason(season string) ([]*Anime, error)
	FindByYearAndSeason(year int, season string) ([]*Anime, error)
	FindByYearAndSeasonAndType(year int, season string, recordType RECORD_TYPE) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatus(year int, season string, recordType RECORD_TYPE, status string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSource(year int, season string, recordType RECORD_TYPE, status string, source string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenre(year int, season string, recordType RECORD_TYPE, status string, source string, genre string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudio(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensors(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRating(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRatingAndName(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string, name string) ([]*Anime, error)
	TopRatedAnime(limit int) ([]*Anime, error)
	MostPopularAnime(limit int) ([]*Anime, error)
	NewestAnime(limit int) ([]*Anime, error)
}

type AnimeRepository struct {
	db *db.DB
}

func NewAnimeRepository(db *db.DB) AnimeRepositoryImpl {
	return &AnimeRepository{db: db}
}

func (a *AnimeRepository) FindAll() ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindById(id string) (*Anime, error) {
	var anime Anime
	err := a.db.DB.Where("id = ?", id).First(&anime).Error
	if err != nil {
		return nil, err
	}
	return &anime, nil
}

func (a *AnimeRepository) FindByName(name string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("name = ?", name).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByType(recordType RECORD_TYPE) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("type = ?", recordType).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByStatus(status string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("status = ?", status).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindBySource(source string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("source = ?", source).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByGenre(genre string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("genre = ?", genre).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByStudio(studio string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("studio = ?", studio).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByLicensors(licensors string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("licensors = ?", licensors).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByRating(rating string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("rating = ?", rating).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYear(year int) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ?", year).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindBySeason(season string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("season = ?", season).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeason(year int, season string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ?", year, season).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndType(year int, season string, recordType RECORD_TYPE) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ?", year, season, recordType).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatus(year int, season string, recordType RECORD_TYPE, status string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ?", year, season, recordType, status).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSource(year int, season string, recordType RECORD_TYPE, status string, source string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ?", year, season, recordType, status, source).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenre(year int, season string, recordType RECORD_TYPE, status string, source string, genre string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ?", year, season, recordType, status, source, genre).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudio(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ?", year, season, recordType, status, source, genre, studio).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensors(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ?", year, season, recordType, status, source, genre, studio, licensors).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRating(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ? AND rating = ?", year, season, recordType, status, source, genre, studio, licensors, rating).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRatingAndName(year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string, name string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ? AND rating = ? AND name = ?", year, season, recordType, status, source, genre, studio, licensors, rating, name).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) TopRatedAnime(limit int) ([]*Anime, error) {
	var animes []*Anime
	// order by rating desc and rating does not equal N/A
	err := a.db.DB.Where("rating != ?", "N/A").Order("rating desc").Limit(limit).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) MostPopularAnime(limit int) ([]*Anime, error) {
	var animes []*Anime
	// order by popularity desc and popularity does not equal N/A
	err := a.db.DB.Where("ranking != ?", "N/A").Order("ranking asc").Limit(limit).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) NewestAnime(limit int) ([]*Anime, error) {
	var animes []*Anime
	// order by start date desc where not null
	err := a.db.DB.Where("start_date ").Order("start_date desc").Limit(limit).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}
