package anime

import (
	"context"
	metrics_lib "github.com/TempMee/go-metrics-lib"
	"github.com/weeb-vip/anime-api/internal/db"
	anime "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"github.com/weeb-vip/anime-api/metrics"
	"time"
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
	AiringAnime(ctx context.Context) ([]*AnimeWithNextEpisode, error)
	AiringAnimeDays(ctx context.Context, startDate *time.Time, days *int) ([]*AnimeWithNextEpisode, error)
	AiringAnimeEndDate(ctx context.Context, startDate *time.Time, endDate *time.Time) ([]*AnimeWithNextEpisode, error)
	SearchAnime(ctx context.Context, search string, page int, limit int) ([]*Anime, error)
}

type AnimeRepository struct {
	db *db.DB
}

func NewAnimeRepository(db *db.DB) AnimeRepositoryImpl {
	return &AnimeRepository{db: db}
}

func (a *AnimeRepository) FindAll(ctx context.Context) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindById(ctx context.Context, id string) (*Anime, error) {
	startTime := time.Now()

	var anime Anime
	err := a.db.DB.Where("id = ?", id).First(&anime).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return &anime, nil
}

func (a *AnimeRepository) FindByName(ctx context.Context, name string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("name = ?", name).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByType(ctx context.Context, recordType RECORD_TYPE) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("type = ?", recordType).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByStatus(ctx context.Context, status string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("status = ?", status).Find(&animes).Error
	if err != nil {

		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindBySource(ctx context.Context, source string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("source = ?", source).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByGenre(ctx context.Context, genre string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("genre = ?", genre).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByStudio(ctx context.Context, studio string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("studio = ?", studio).Find(&animes).Error
	if err != nil {

		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByLicensors(ctx context.Context, licensors string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("licensors = ?", licensors).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByRating(ctx context.Context, rating string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("rating = ?", rating).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYear(ctx context.Context, year int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ?", year).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindBySeason(ctx context.Context, season string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("season = ?", season).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeason(ctx context.Context, year int, season string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ?", year, season).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndType(ctx context.Context, year int, season string, recordType RECORD_TYPE) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ?", year, season, recordType).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatus(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ?", year, season, recordType, status).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSource(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ?", year, season, recordType, status, source).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenre(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ?", year, season, recordType, status, source, genre).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudio(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ?", year, season, recordType, status, source, genre, studio).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensors(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ?", year, season, recordType, status, source, genre, studio, licensors).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRating(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ? AND rating = ?", year, season, recordType, status, source, genre, studio, licensors, rating).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) FindByYearAndSeasonAndTypeAndStatusAndSourceAndGenreAndStudioAndLicensorsAndRatingAndName(ctx context.Context, year int, season string, recordType RECORD_TYPE, status string, source string, genre string, studio string, licensors string, rating string, name string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("year = ? AND season = ? AND type = ? AND status = ? AND source = ? AND genre = ? AND studio = ? AND licensors = ? AND rating = ? AND name = ?", year, season, recordType, status, source, genre, studio, licensors, rating, name).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) TopRatedAnime(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	// order by rating desc and rating does not equal N/A
	err := a.db.DB.Where("rating != ?", "N/A").Order("rating desc").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) MostPopularAnime(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	// order by popularity desc and popularity does not equal N/A
	err := a.db.DB.Where("ranking != ?", "N/A").Order("ranking asc").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) NewestAnime(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	// order by start date desc where not null
	err := a.db.DB.Where("created_at ").Order("created_at desc").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) AiringAnime(ctx context.Context) ([]*AnimeWithNextEpisode, error) {
	startTime := time.Now()

	var animes []*AnimeWithNextEpisode

	subQuery := a.db.DB.Model(&anime.AnimeEpisode{}).
		Select("anime_id, MIN(aired) AS next_aired").
		// curr date inclusive
		Where("aired BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 30 DAY)").
		Group("anime_id")

	err := a.db.DB.Table("anime").
		Select("anime.*"). // Only scan anime fields into AnimeWithNextEpisode
		Joins("JOIN (?) AS e ON anime.id = e.anime_id", subQuery).
		Where("anime.end_date IS NULL").
		Order("e.next_aired").
		Scan(&animes).Error

	metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  map[bool]string{true: metrics_lib.Success, false: metrics_lib.Error}[err == nil],
	})

	if err != nil {
		return nil, err
	}

	// Populate the next_aired field in each AnimeWithNextEpisode
	for i := range animes {
		var nextEpisode anime.AnimeEpisode
		err := a.db.DB.Model(&anime.AnimeEpisode{}).
			Where("anime_id = ? AND aired BETWEEN CURDATE() AND DATE_ADD(CURDATE(), INTERVAL 30 DAY)", animes[i].ID).
			Order("aired").
			First(&nextEpisode).Error
		if err != nil {
			return nil, err
		}
		animes[i].NextEpisode = &nextEpisode
	}
	return animes, nil
}

func (a *AnimeRepository) SearchAnime(ctx context.Context, search string, page int, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.Where("title_en LIKE ? OR title_jp LIKE ? OR title_synonyms LIKE ? OR title_romaji LIKE ? OR title_kanji LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Limit(limit).Offset((page - 1) * limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) AiringAnimeDays(ctx context.Context, startDate *time.Time, days *int) ([]*AnimeWithNextEpisode, error) {
	startTime := time.Now()

	var animes []*AnimeWithNextEpisode

	subQuery := a.db.DB.Model(&anime.AnimeEpisode{}).
		Select("anime_id, MIN(aired) AS next_aired").
		Where("aired BETWEEN ? AND DATE_ADD(?, INTERVAL ? DAY)", startDate, startDate, days).
		Group("anime_id")

	err := a.db.DB.Table("anime").
		Select("anime.*"). // Only scan anime fields into AnimeWithNextEpisode
		Joins("JOIN (?) AS e ON anime.id = e.anime_id", subQuery).
		Order("e.next_aired").
		Scan(&animes).Error

	if err != nil {
		return nil, err
	}

	// Populate the next_aired field in each AnimeWithNextEpisode
	for i := range animes {
		var nextEpisode anime.AnimeEpisode
		err := a.db.DB.Model(&anime.AnimeEpisode{}).
			Where("anime_id = ? AND aired BETWEEN ? AND DATE_ADD(?, INTERVAL ? DAY)", animes[i].ID, startDate, startDate, days).
			Order("aired").
			First(&nextEpisode).Error
		if err != nil {
			return nil, err
		}
		animes[i].NextEpisode = &nextEpisode
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})
	return animes, nil
}

func (a *AnimeRepository) AiringAnimeEndDate(ctx context.Context, startDate *time.Time, endDate *time.Time) ([]*AnimeWithNextEpisode, error) {
	startTime := time.Now()

	var animes []*AnimeWithNextEpisode

	subQuery := a.db.DB.Model(&anime.AnimeEpisode{}).
		Select("anime_id, MIN(aired) AS next_aired").
		Where("aired BETWEEN ? AND ?", startDate, endDate).
		Group("anime_id")

	err := a.db.DB.Table("anime").
		Select("anime.*"). // Only scan anime fields into AnimeWithNextEpisode
		Joins("JOIN (?) AS e ON anime.id = e.anime_id", subQuery).
		Order("e.next_aired").
		Scan(&animes).Error

	if err != nil {
		return nil, err
	}

	// Populate the next_aired field in each AnimeWithNextEpisode
	for i := range animes {
		var nextEpisode anime.AnimeEpisode
		err := a.db.DB.Model(&anime.AnimeEpisode{}).
			Where("anime_id = ? AND aired BETWEEN ? AND ?", animes[i].ID, startDate, endDate).
			Order("aired").
			First(&nextEpisode).Error
		if err != nil {
			return nil, err
		}
		animes[i].NextEpisode = &nextEpisode
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
	})

	return animes, nil
}
