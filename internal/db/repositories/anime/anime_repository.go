package anime

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type RECORD_TYPE string

type AnimeRepositoryImpl interface {
	FindAll(ctx context.Context) ([]*Anime, error)
	FindById(ctx context.Context, id string) (*Anime, error)
	FindByName(ctx context.Context, name string) ([]*Anime, error)
	FindByType(ctx context.Context, recordType RECORD_TYPE) ([]*Anime, error)
	FindByStatus(ctx context.Context, status string) ([]*Anime, error)
	FindBySource(ctx context.Context, source string) ([]*Anime, error)
	FindByGenre(ctx context.Context, genre string) ([]*Anime, error)
	FindByStudio(ctx context.Context, studio string) ([]*Anime, error)
	FindByLicensors(ctx context.Context, licensors string) ([]*Anime, error)
	FindByRating(ctx context.Context, rating string) ([]*Anime, error)
	FindByYear(ctx context.Context, year int) ([]*Anime, error)
	FindBySeason(ctx context.Context, season string) ([]*Anime, error)
	FindByYearAndSeason(ctx context.Context, year int, season string) ([]*Anime, error)
	FindByYearAndSeasonAndType(ctx context.Context, year int, season string, recordType RECORD_TYPE) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatus(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSource(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenre(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudio(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensors(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRating(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string) ([]*Anime, error)
	FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRatingAndName(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string, name string) ([]*Anime, error)
	TopRatedAnime(ctx context.Context, limit int) ([]*Anime, error)
	MostPopularAnime(ctx context.Context, limit int) ([]*Anime, error)
	NewestAnime(ctx context.Context, limit int) ([]*Anime, error)
}

type AnimeRepository struct {
	db *db.DB
}

func NewAnimeRepository(db *db.DB) AnimeRepositoryImpl {
	return &AnimeRepository{db: db}
}

func (a *AnimeRepository) FindAll(ctx context.Context) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindById(ctx context.Context, id string) (*Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "FindById")
	span.SetTag("type", "repository")
	defer span.Finish()

	var anime Anime
	err := a.db.DB.Where("id = ?", id).First(&anime).Error
	if err != nil {
		return nil, err
	}
	return &anime, nil
}

func (a *AnimeRepository) FindByName(ctx context.Context, name string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("name = ?", name).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByType(ctx context.Context, recordType RECORD_TYPE) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("type = ?", recordType).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByStatus(ctx context.Context, status string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("status = ?", status).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindBySource(ctx context.Context, source string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("source = ?", source).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByGenre(ctx context.Context, genre string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("genre = ?", genre).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByStudio(ctx context.Context, studio string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("studio = ?", studio).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByLicensors(ctx context.Context, licensors string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("licensors = ?", licensors).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByRating(ctx context.Context, rating string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("rating = ?", rating).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYear(ctx context.Context, year int) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ?", year).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindBySeason(ctx context.Context, season string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("season = ?", season).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeason(ctx context.Context, year int, season string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ?", year, season).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndType(ctx context.Context, year int, season string, recordType RECORD_TYPE) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ?", year, season, recordType).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatus(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ?", year, season, recordType, status).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSource(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ?", year, season, recordType, status, source).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenre(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ?", year, season, recordType, status, source, genre).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudio(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ?", year, season, recordType, status, source, genre, studio).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensors(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ?", year, season, recordType, status, source, genre, studio, licensors).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRating(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ? AND rating = ?", year, season, recordType, status, source, genre, studio, licensors, rating).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRatingAndName(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string, name string) ([]*Anime, error) {
	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ? AND rating = ? AND name = ?", year, season, recordType, status, source, genre, studio, licensors, rating, name).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) TopRatedAnime(ctx context.Context, limit int) ([]*Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "TopRatedAnime")
	span.SetTag("type", "repository")
	defer span.Finish()

	var animes []*Anime
	// order by rating desc and rating does not equal N/A
	err := a.db.DB.Where("rating != ?", "N/A").Order("rating desc").Limit(limit).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) MostPopularAnime(ctx context.Context, limit int) ([]*Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "MostPopularAnime")
	span.SetTag("type", "repository")
	defer span.Finish()

	var animes []*Anime
	// order by popularity desc and popularity does not equal N/A
	err := a.db.DB.Where("ranking != ?", "N/A").Order("ranking asc").Limit(limit).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}

func (a *AnimeRepository) NewestAnime(ctx context.Context, limit int) ([]*Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "NewestAnime")
	span.SetTag("type", "repository")
	defer span.Finish()

	var animes []*Anime
	// order by start date desc where not null
	err := a.db.DB.Where("start_date ").Order("start_date desc").Limit(limit).Find(&animes).Error
	if err != nil {
		return nil, err
	}
	return animes, nil
}
