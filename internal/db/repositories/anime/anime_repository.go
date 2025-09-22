package anime

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db"
	anime "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
	"gorm.io/gorm"
	"time"
)

type RECORD_TYPE string

var tzTokyo = mustLoadTZ("Asia/Tokyo")

func mustLoadTZ(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

type AnimeRepositoryImpl interface {
	FindAll(ctx context.Context) ([]*Anime, error)
	FindAllWithEpisodes(ctx context.Context) ([]*Anime, error)
	FindById(ctx context.Context, id string) (*Anime, error)
	FindByIdWithEpisodes(ctx context.Context, id string) (*Anime, error)
	FindByIDs(ctx context.Context, ids []string) ([]*Anime, error)
	FindByIDsWithEpisodes(ctx context.Context, ids []string) ([]*Anime, error)
	FindByName(ctx context.Context, name string) ([]*Anime, error)
	FindByNameWithEpisodes(ctx context.Context, name string) ([]*Anime, error)
	FindByType(ctx context.Context, recordType RECORD_TYPE) ([]*Anime, error)
	FindByTypeWithEpisodes(ctx context.Context, recordType RECORD_TYPE) ([]*Anime, error)
	FindByStatus(ctx context.Context, status string) ([]*Anime, error)
	FindByStatusWithEpisodes(ctx context.Context, status string) ([]*Anime, error)
	FindBySource(ctx context.Context, source string) ([]*Anime, error)
	FindByGenre(ctx context.Context, genre string) ([]*Anime, error)
	FindByStudio(ctx context.Context, studio string) ([]*Anime, error)
	FindByLicensors(ctx context.Context, licensors string) ([]*Anime, error)
	FindByRating(ctx context.Context, rating string) ([]*Anime, error)
	FindByYear(ctx context.Context, year int) ([]*Anime, error)
	TopRatedAnime(ctx context.Context, limit int) ([]*Anime, error)
	TopRatedAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error)
	MostPopularAnime(ctx context.Context, limit int) ([]*Anime, error)
	MostPopularAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error)
	NewestAnime(ctx context.Context, limit int) ([]*Anime, error)
	NewestAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error)
	AiringAnime(ctx context.Context) ([]*AnimeWithNextEpisode, error)
	AiringAnimeDays(ctx context.Context, startDate *time.Time, days *int) ([]*AnimeWithNextEpisode, error)
	AiringAnimeEndDate(ctx context.Context, startDate *time.Time, endDate *time.Time) ([]*AnimeWithNextEpisode, error)
	AiringAnimeWithEpisodes(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*Anime, error)
	SearchAnime(ctx context.Context, search string, page int, limit int) ([]*Anime, error)
	SearchAnimeWithEpisodes(ctx context.Context, search string, page int, limit int) ([]*Anime, error)
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
	err := a.db.DB.WithContext(ctx).Find(&animes).Error
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

func (a *AnimeRepository) FindAllWithEpisodes(ctx context.Context) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("id = ?", id).First(&anime).Error
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

func (a *AnimeRepository) FindByIdWithEpisodes(ctx context.Context, id string) (*Anime, error) {
	startTime := time.Now()

	var anime Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("id = ?", id).First(&anime).Error
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
	err := a.db.DB.WithContext(ctx).Where("name = ?", name).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("type = ?", recordType).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("status = ?", status).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("source = ?", source).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("genre = ?", genre).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("studio = ?", studio).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("licensors = ?", licensors).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("rating = ?", rating).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("year = ?", year).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("rating != ?", "N/A").Order("rating desc").Limit(limit).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("ranking != ?", "N/A").Order("ranking asc").Limit(limit).Find(&animes).Error
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
	err := a.db.DB.WithContext(ctx).Where("created_at ").Order("created_at desc").Limit(limit).Find(&animes).Error
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

	nowJST := startOfDayIn(time.Now().UTC(), tzTokyo)
	endJST := nowJST.AddDate(0, 0, 30)

	subQuery := a.db.DB.WithContext(ctx).Model(&anime.AnimeEpisode{}).
		Select("anime_id, MIN(aired) AS next_aired").
		Where("aired BETWEEN ? AND ?", nowJST, endJST).
		Group("anime_id")

	err := a.db.DB.WithContext(ctx).Table("anime").
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
		err := a.db.DB.WithContext(ctx).Model(&anime.AnimeEpisode{}).
			Where("anime_id = ? AND aired BETWEEN ? AND ?", animes[i].ID, nowJST, endJST).
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
	err := a.db.DB.WithContext(ctx).Where("title_en LIKE ? OR title_jp LIKE ? OR title_synonyms LIKE ? OR title_romaji LIKE ? OR title_kanji LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Limit(limit).Offset((page - 1) * limit).Find(&animes).Error
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

	startJST := startOfDayIn(startDate.UTC(), tzTokyo)
	endJST := startJST.AddDate(0, 0, *days)

	var animes []*AnimeWithNextEpisode

	subQuery := a.db.DB.WithContext(ctx).Model(&anime.AnimeEpisode{}).
		Select("anime_id, MIN(aired) AS next_aired").
		Where("aired BETWEEN ? AND ?", startJST, endJST).
		Group("anime_id")

	err := a.db.DB.WithContext(ctx).Table("anime").
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
		err := a.db.DB.WithContext(ctx).Model(&anime.AnimeEpisode{}).
			Where("anime_id = ? AND aired BETWEEN ? AND ?", animes[i].ID, startJST, endJST).
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

	startJST := startOfDayIn(startDate.UTC(), tzTokyo)
	// keep end as provided but in JST zone (no implicit +1 day)
	endJST := endDate.In(tzTokyo)

	subQuery := a.db.DB.WithContext(ctx).Model(&anime.AnimeEpisode{}).
		Select("anime_id, MIN(aired) AS next_aired").
		Where("aired BETWEEN ? AND ?", startJST, endJST).
		Group("anime_id")

	err := a.db.DB.WithContext(ctx).Table("anime").
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
		err := a.db.DB.WithContext(ctx).Model(&anime.AnimeEpisode{}).
			Where("anime_id = ? AND aired BETWEEN ? AND ?", animes[i].ID, startJST, endJST).
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

func (a *AnimeRepository) FindByNameWithEpisodes(ctx context.Context, name string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("name = ?", name).Find(&animes).Error
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

func (a *AnimeRepository) FindByTypeWithEpisodes(ctx context.Context, recordType RECORD_TYPE) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("type = ?", recordType).Find(&animes).Error
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

func (a *AnimeRepository) FindByStatusWithEpisodes(ctx context.Context, status string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("status = ?", status).Find(&animes).Error
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

func (a *AnimeRepository) TopRatedAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("rating != ?", "N/A").Order("rating desc").Limit(limit).Find(&animes).Error
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

func (a *AnimeRepository) MostPopularAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("ranking != ?", "N/A").Order("ranking asc").Limit(limit).Find(&animes).Error
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

func (a *AnimeRepository) NewestAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("created_at ").Order("created_at desc").Limit(limit).Find(&animes).Error
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

func (a *AnimeRepository) SearchAnimeWithEpisodes(ctx context.Context, search string, page int, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("title_en LIKE ? OR title_jp LIKE ? OR title_synonyms LIKE ? OR title_romaji LIKE ? OR title_kanji LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%").Limit(limit).Offset((page - 1) * limit).Find(&animes).Error
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

func (a *AnimeRepository) AiringAnimeWithEpisodes(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime

	query := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	})

	// Handle date filtering logic similar to the service layer
	if startDate != nil && endDate != nil {
		// Filter anime that have episodes airing between startDate and endDate
		query = query.Joins("INNER JOIN episodes ON anime.id = episodes.anime_id").
			Where("episodes.aired BETWEEN ? AND ?", *startDate, *endDate).
			Where("anime.end_date IS NULL OR anime.end_date >= ?", startDate)
	} else if startDate != nil && days != nil {
		// Filter anime that have episodes airing from startDate for the next N days
		endTime := startDate.AddDate(0, 0, *days)
		query = query.Joins("INNER JOIN episodes ON anime.id = episodes.anime_id").
			Where("episodes.aired BETWEEN ? AND ?", *startDate, endTime).
			Where("anime.end_date IS NULL OR anime.end_date >= ?", startDate)
	} else {
		// Default: currently airing anime (next 30 days)
		nowJST := startOfDayIn(time.Now().UTC(), tzTokyo)
		endJST := nowJST.AddDate(0, 0, 30)
		query = query.Joins("INNER JOIN episodes ON anime.id = episodes.anime_id").
			Where("episodes.aired BETWEEN ? AND ?", nowJST, endJST).
			Where("anime.end_date IS NULL")
	}

	err := query.Distinct("anime.*").Find(&animes).Error

	metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  map[bool]string{true: metrics_lib.Success, false: metrics_lib.Error}[err == nil],
	})

	if err != nil {
		return nil, err
	}

	return animes, nil
}

// Start-of-day in a specific zone (keeps date math stable)
func startOfDayIn(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}
