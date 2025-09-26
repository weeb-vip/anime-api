package anime

import (
	"context"
	"fmt"
	"time"

	"github.com/weeb-vip/anime-api/internal/cache"
	"github.com/weeb-vip/anime-api/internal/db"
	animeEpisode "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"github.com/weeb-vip/anime-api/internal/logger"
	"github.com/weeb-vip/anime-api/metrics"
	"github.com/weeb-vip/anime-api/tracing"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
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
	FindBySeasonWithEpisodes(ctx context.Context, season string) ([]*Anime, error)
	FindBySeasonWithEpisodesOptimized(ctx context.Context, season string) ([]*Anime, error)
	FindBySeasonWithIndexHints(ctx context.Context, season string) ([]*Anime, error)
	FindBySeasonBatched(ctx context.Context, season string) ([]*Anime, error)
	FindBySeasonAnimeOnlyOptimized(ctx context.Context, season string) ([]*Anime, error)
	FindBySeasonWithFieldSelection(ctx context.Context, season string, fields *FieldSelection, limit int) ([]*Anime, error)
}

type AnimeRepository struct {
	db    *db.DB
	cache *cache.CacheService
}

func NewAnimeRepository(db *db.DB) AnimeRepositoryImpl {
	return &AnimeRepository{db: db, cache: nil}
}

func NewAnimeRepositoryWithCache(db *db.DB, cacheService *cache.CacheService) AnimeRepositoryImpl {
	return &AnimeRepository{db: db, cache: cacheService}
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) FindById(ctx context.Context, id string) (*Anime, error) {
	// Try cache first if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().AnimeByID(id)
		var anime Anime
		err := a.cache.GetJSON(ctx, key, &anime)
		if err == nil {
			return &anime, nil
		}
		// Continue to database if cache miss or error
	}

	startTime := time.Now()
	var anime Anime
	err := a.db.DB.WithContext(ctx).Where("id = ?", id).First(&anime).Error
	if err != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			"anime",
			"select",
			metrics.Error,
		)
		return nil, err
	}

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		"anime",
		"select",
		metrics.Success,
	)

	// Store in cache if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().AnimeByID(id)
		_ = a.cache.SetJSON(ctx, key, &anime, a.cache.GetAnimeDataTTL())
	}

	return &anime, nil
}

func (a *AnimeRepository) FindByIdWithEpisodes(ctx context.Context, id string) (*Anime, error) {
	// Try cache first if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().AnimeWithEpisodesByID(id)
		var anime Anime
		err := a.cache.GetJSON(ctx, key, &anime)
		if err == nil {
			return &anime, nil
		}
		// Continue to database if cache miss or error
	}

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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})

	// Store in cache if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().AnimeWithEpisodesByID(id)
		_ = a.cache.SetJSON(ctx, key, &anime, a.cache.GetAnimeDataTTL())
	}

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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) TopRatedAnime(ctx context.Context, limit int) ([]*Anime, error) {
	// Try cache first if available
	if a.cache != nil {
		key := fmt.Sprintf("%s:top_rated:%d", a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1], limit)
		var animeList []*Anime
		err := a.cache.GetJSON(ctx, key, &animeList)
		if err == nil {
			return animeList, nil
		}
		// Continue to database if cache miss or error
	}

	startTime := time.Now()

	var animes []*Anime
	// Use numeric rating for better performance - automatically excludes NULL ratings (N/A)
	err := a.db.DB.WithContext(ctx).Where("rating IS NOT NULL").Order("rating desc, id").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})

	// Store in cache if available
	if a.cache != nil {
		key := fmt.Sprintf("%s:top_rated:%d", a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1], limit)
		_ = a.cache.SetJSON(ctx, key, animes, a.cache.GetAnimeDataTTL())
	}

	return animes, nil
}

func (a *AnimeRepository) MostPopularAnime(ctx context.Context, limit int) ([]*Anime, error) {
	// Try cache first if available
	if a.cache != nil {
		key := fmt.Sprintf("%s:most_popular:%d", a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1], limit)
		var animeList []*Anime
		err := a.cache.GetJSON(ctx, key, &animeList)
		if err == nil {
			return animeList, nil
		}
		// Continue to database if cache miss or error
	}

	startTime := time.Now()

	var animes []*Anime
	// Order by ranking asc (lower ranking = more popular) - only include anime with rankings
	err := a.db.DB.WithContext(ctx).Where("ranking IS NOT NULL").Order("ranking asc, id").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) NewestAnime(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	// Order by created_at desc (newest first) - all anime have created_at so no WHERE needed
	err := a.db.DB.WithContext(ctx).Order("created_at desc, id").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) AiringAnime(ctx context.Context) ([]*AnimeWithNextEpisode, error) {
	startTime := time.Now()

	var animes []*AnimeWithNextEpisode

	nowJST := startOfDayIn(time.Now().UTC(), tzTokyo)
	endJST := nowJST.AddDate(0, 0, 30)

	subQuery := a.db.DB.WithContext(ctx).Model(&animeEpisode.AnimeEpisode{}).
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
		Env:     metrics.GetCurrentEnv(),
	})

	if err != nil {
		return nil, err
	}

	// Populate the next_aired field in each AnimeWithNextEpisode
	for i := range animes {
		var nextEpisode animeEpisode.AnimeEpisode
		err := a.db.DB.WithContext(ctx).Model(&animeEpisode.AnimeEpisode{}).
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) AiringAnimeDays(ctx context.Context, startDate *time.Time, days *int) ([]*AnimeWithNextEpisode, error) {
	startTime := time.Now()

	startJST := startOfDayIn(startDate.UTC(), tzTokyo)
	endJST := startJST.AddDate(0, 0, *days)

	var animes []*AnimeWithNextEpisode

	subQuery := a.db.DB.WithContext(ctx).Model(&animeEpisode.AnimeEpisode{}).
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
		var nextEpisode animeEpisode.AnimeEpisode
		err := a.db.DB.WithContext(ctx).Model(&animeEpisode.AnimeEpisode{}).
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
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) AiringAnimeEndDate(ctx context.Context, startDate *time.Time, endDate *time.Time) ([]*AnimeWithNextEpisode, error) {
	startTime := time.Now()

	var animes []*AnimeWithNextEpisode

	startJST := startOfDayIn(startDate.UTC(), tzTokyo)
	// keep end as provided but in JST zone (no implicit +1 day)
	endJST := endDate.In(tzTokyo)

	subQuery := a.db.DB.WithContext(ctx).Model(&animeEpisode.AnimeEpisode{}).
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
		var nextEpisode animeEpisode.AnimeEpisode
		err := a.db.DB.WithContext(ctx).Model(&animeEpisode.AnimeEpisode{}).
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
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) TopRatedAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("rating IS NOT NULL").Order("rating desc, id").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) MostPopularAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Where("ranking IS NOT NULL").Order("ranking asc, id").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) NewestAnimeWithEpisodes(ctx context.Context, limit int) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	}).Order("created_at desc, id").Limit(limit).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

func (a *AnimeRepository) SearchAnimeWithEpisodes(ctx context.Context, search string, page int, limit int) ([]*Anime, error) {
	// Try cache first if available
	if a.cache != nil {
		key := fmt.Sprintf("%s:search_episodes:%s:page_%d:limit_%d", a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1], search, page, limit)
		var animeList []*Anime
		err := a.cache.GetJSON(ctx, key, &animeList)
		if err == nil {
			return animeList, nil
		}
		// Continue to database if cache miss or error
	}

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
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})

	// Store in cache if available
	if a.cache != nil {
		key := fmt.Sprintf("%s:search_episodes:%s:page_%d:limit_%d", a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1], search, page, limit)
		_ = a.cache.SetJSON(ctx, key, animes, a.cache.GetAnimeDataTTL())
	}

	return animes, nil
}

func (a *AnimeRepository) AiringAnimeWithEpisodes(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*Anime, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeRepository.AiringAnimeWithEpisodes",
		trace.WithAttributes(
			attribute.String("repository", "anime"),
			attribute.String("method", "AiringAnimeWithEpisodes"),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	// Generate cache key based on parameters
	cacheKey := ""
	if a.cache != nil {
		// Add cache key generation tracing
		_, cacheKeySpan := tracer.Start(ctx, "AnimeRepository.GenerateCacheKey",
			trace.WithAttributes(
				attribute.String("cache.operation", "key_generation"),
			),
			trace.WithSpanKind(trace.SpanKindInternal),
		)

		if startDate != nil && endDate != nil {
			cacheKey = fmt.Sprintf("%s:airing:start_%s:end_%s",
				a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1],
				startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
		} else if startDate != nil && days != nil {
			cacheKey = fmt.Sprintf("%s:airing:start_%s:days_%d",
				a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1],
				startDate.Format("2006-01-02"), *days)
		} else {
			// Default case (next 30 days from now)
			cacheKey = fmt.Sprintf("%s:airing:default",
				a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1])
		}

		cacheKeySpan.SetAttributes(attribute.String("cache.key", cacheKey))
		cacheKeySpan.End()

		// Try cache first with tracing
		_, cacheGetSpan := tracer.Start(ctx, "AnimeRepository.CacheGet",
			trace.WithAttributes(
				attribute.String("cache.operation", "get"),
				attribute.String("cache.key", cacheKey),
			),
			trace.WithSpanKind(trace.SpanKindInternal),
		)

		var animeList []*Anime
		err := a.cache.GetJSON(ctx, cacheKey, &animeList)
		if err == nil {
			cacheGetSpan.SetAttributes(
				attribute.String("cache.result", "hit"),
				attribute.Int("cache.items_count", len(animeList)),
			)
			cacheGetSpan.SetStatus(codes.Ok, "cache hit")
			cacheGetSpan.End()
			span.SetAttributes(attribute.String("data_source", "cache"))

			// Debug logging for cache hits
			log := logger.FromCtx(ctx)
			log.Info().
				Str("cache_key", cacheKey).
				Int("items_count", len(animeList)).
				Msg("Cache HIT for AiringAnimeWithEpisodes")

			return animeList, nil
		}

		cacheGetSpan.SetAttributes(attribute.String("cache.result", "miss"))
		cacheGetSpan.RecordError(err)
		cacheGetSpan.SetStatus(codes.Error, "cache miss: "+err.Error())
		cacheGetSpan.End()

		// Debug logging for cache misses
		log := logger.FromCtx(ctx)
		log.Info().
			Str("cache_key", cacheKey).
			Err(err).
			Msg("Cache MISS for AiringAnimeWithEpisodes")

		// Continue to database if cache miss or error
	}

	// Add database query tracing
	_, dbQuerySpan := tracer.Start(ctx, "AnimeRepository.DatabaseQuery",
		trace.WithAttributes(
			attribute.String("db.operation", "select_with_episodes"),
			attribute.String("query.type", "airing_anime"),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
	)

	startTime := time.Now()
	var animes []*Anime

	query := a.db.DB.WithContext(ctx).Preload("AnimeEpisodes", func(db *gorm.DB) *gorm.DB {
		return db.Order("episode ASC")
	})

	// Use optimized subquery approach to avoid DISTINCT and improve performance
	var subquery *gorm.DB

	if startDate != nil && endDate != nil {
		// Subquery to get unique anime IDs with episodes in date range
		subquery = a.db.DB.Model(&animeEpisode.AnimeEpisode{}).
			Select("DISTINCT anime_id").
			Where("aired BETWEEN ? AND ?", *startDate, *endDate)

		// Filter anime by subquery results and end_date condition
		query = query.Where("anime.id IN (?)", subquery).
			Where("anime.end_date IS NULL OR anime.end_date >= ?", startDate)

	} else if startDate != nil && days != nil {
		// Subquery for days range
		endTime := startDate.AddDate(0, 0, *days)
		subquery = a.db.DB.Model(&animeEpisode.AnimeEpisode{}).
			Select("DISTINCT anime_id").
			Where("aired BETWEEN ? AND ?", *startDate, endTime)

		// Filter anime by subquery results and end_date condition
		query = query.Where("anime.id IN (?)", subquery).
			Where("anime.end_date IS NULL OR anime.end_date >= ?", startDate)

	} else {
		// Default: currently airing anime (next 30 days)
		nowJST := startOfDayIn(time.Now().UTC(), tzTokyo)
		endJST := nowJST.AddDate(0, 0, 30)
		subquery = a.db.DB.Model(&animeEpisode.AnimeEpisode{}).
			Select("DISTINCT anime_id").
			Where("aired BETWEEN ? AND ?", nowJST, endJST)

		// Filter anime by subquery results (no end_date check for default case)
		query = query.Where("anime.id IN (?)", subquery).
			Where("anime.end_date IS NULL")
	}

	err := query.Find(&animes).Error

	metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  map[bool]string{true: metrics_lib.Success, false: metrics_lib.Error}[err == nil],
		Env:     metrics.GetCurrentEnv(),
	})

	// Complete database query tracing
	dbQuerySpan.SetAttributes(
		attribute.Int("db.rows_returned", len(animes)),
		attribute.Int64("db.duration_ms", time.Since(startTime).Milliseconds()),
	)

	if err != nil {
		dbQuerySpan.RecordError(err)
		dbQuerySpan.SetStatus(codes.Error, err.Error())
		dbQuerySpan.End()
		return nil, err
	}

	dbQuerySpan.SetStatus(codes.Ok, "query completed")
	dbQuerySpan.End()
	span.SetAttributes(attribute.String("data_source", "database"))

	// Store in cache if available with tracing
	if a.cache != nil && cacheKey != "" {
		_, cacheSetSpan := tracer.Start(ctx, "AnimeRepository.CacheSet",
			trace.WithAttributes(
				attribute.String("cache.operation", "set"),
				attribute.String("cache.key", cacheKey),
				attribute.Int("cache.items_count", len(animes)),
			),
			trace.WithSpanKind(trace.SpanKindInternal),
		)

		err := a.cache.SetJSON(ctx, cacheKey, animes, a.cache.GetAnimeDataTTL())
		if err != nil {
			cacheSetSpan.RecordError(err)
			cacheSetSpan.SetStatus(codes.Error, "cache set failed: "+err.Error())

			// Debug logging for cache set failures
			log := logger.FromCtx(ctx)
			log.Error().
				Str("cache_key", cacheKey).
				Int("items_count", len(animes)).
				Err(err).
				Msg("Failed to SET cache for AiringAnimeWithEpisodes")
		} else {
			cacheSetSpan.SetStatus(codes.Ok, "cache set successful")

			// Debug logging for successful cache sets
			log := logger.FromCtx(ctx)
			log.Info().
				Str("cache_key", cacheKey).
				Int("items_count", len(animes)).
				Msg("Successfully SET cache for AiringAnimeWithEpisodes")
		}
		cacheSetSpan.End()
	}

	return animes, nil
}

// FindBySeasonWithEpisodes performs a single optimized query with joins to get anime and episodes by season
func (a *AnimeRepository) FindBySeasonWithEpisodes(ctx context.Context, season string) ([]*Anime, error) {
	startTime := time.Now()

	// Create a custom query that joins anime_seasons -> anime -> episodes in one query
	type AnimeWithEpisodeResult struct {
		// Anime fields
		AnimeID            string     `gorm:"column:anime_id"`
		AnidbID            *string    `gorm:"column:anidbid"`
		TheTVDBID          *string    `gorm:"column:thetvdbid"`
		Type               *RECORD_TYPE `gorm:"column:type"`
		TitleEn            *string    `gorm:"column:title_en"`
		TitleJp            *string    `gorm:"column:title_jp"`
		TitleRomaji        *string    `gorm:"column:title_romaji"`
		TitleKanji         *string    `gorm:"column:title_kanji"`
		TitleSynonyms      *string    `gorm:"column:title_synonyms"`
		ImageURL           *string    `gorm:"column:image_url"`
		Synopsis           *string    `gorm:"column:synopsis"`
		Episodes           *int       `gorm:"column:episodes"`
		Status             *string    `gorm:"column:status"`
		StartDate          *string    `gorm:"column:start_date"`
		EndDate            *string    `gorm:"column:end_date"`
		Genres             *string    `gorm:"column:genres"`
		Duration           *string    `gorm:"column:duration"`
		Broadcast          *string    `gorm:"column:broadcast"`
		Source             *string    `gorm:"column:source"`
		Licensors          *string    `gorm:"column:licensors"`
		Studios            *string    `gorm:"column:studios"`
		Rating             *float64   `gorm:"column:rating"`
		Ranking            *int       `gorm:"column:ranking"`
		AnimeCreatedAt     time.Time  `gorm:"column:anime_created_at"`
		AnimeUpdatedAt     time.Time  `gorm:"column:anime_updated_at"`

		// Episode fields (nullable)
		EpisodeID          *string    `gorm:"column:episode_id"`
		EpisodeNumber      *int       `gorm:"column:episode"`
		EpisodeTitleEn     *string    `gorm:"column:episode_title_en"`
		EpisodeTitleJp     *string    `gorm:"column:episode_title_jp"`
		EpisodeAired       *time.Time `gorm:"column:aired"`
		EpisodeSynopsis    *string    `gorm:"column:episode_synopsis"`
		EpisodeCreatedAt   *time.Time `gorm:"column:episode_created_at"`
		EpisodeUpdatedAt   *time.Time `gorm:"column:episode_updated_at"`
	}

	var results []AnimeWithEpisodeResult
	err := a.db.DB.WithContext(ctx).
		Select(`
			a.id as anime_id,
			a.anidbid,
			a.thetvdbid,
			a.type,
			a.title_en,
			a.title_jp,
			a.title_romaji,
			a.title_kanji,
			a.title_synonyms,
			a.image_url,
			a.synopsis,
			a.episodes,
			a.status,
			a.start_date,
			a.end_date,
			a.genres,
			a.duration,
			a.broadcast,
			a.source,
			a.licensors,
			a.studios,
			a.rating,
			a.ranking,
			a.created_at as anime_created_at,
			a.updated_at as anime_updated_at,
			e.id as episode_id,
			e.episode,
			e.title_en as episode_title_en,
			e.title_jp as episode_title_jp,
			e.aired,
			e.synopsis as episode_synopsis,
			e.created_at as episode_created_at,
			e.updated_at as episode_updated_at
		`).
		Table("anime_seasons as s").
		Joins("INNER JOIN anime as a ON s.anime_id = a.id").
		Joins("LEFT JOIN episodes as e ON a.id = e.anime_id").
		Where("s.season = ?", season).
		Order("a.id ASC, e.episode ASC").
		Find(&results).Error

	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Group results by anime and construct the anime objects with episodes
	animeMap := make(map[string]*Anime)
	for _, result := range results {
		// Create anime if not exists
		if _, exists := animeMap[result.AnimeID]; !exists {
			animeMap[result.AnimeID] = &Anime{
				ID:            result.AnimeID,
				AnidbID:       result.AnidbID,
				TheTVDBID:     result.TheTVDBID,
				Type:          result.Type,
				TitleEn:       result.TitleEn,
				TitleJp:       result.TitleJp,
				TitleRomaji:   result.TitleRomaji,
				TitleKanji:    result.TitleKanji,
				TitleSynonyms: result.TitleSynonyms,
				ImageURL:      result.ImageURL,
				Synopsis:      result.Synopsis,
				Episodes:      result.Episodes,
				Status:        result.Status,
				StartDate:     result.StartDate,
				EndDate:       result.EndDate,
				Genres:        result.Genres,
				Duration:      result.Duration,
				Broadcast:     result.Broadcast,
				Source:        result.Source,
				Licensors:     result.Licensors,
				Studios:       result.Studios,
				Rating:        result.Rating,
				Ranking:       result.Ranking,
				CreatedAt:     result.AnimeCreatedAt,
				UpdatedAt:     result.AnimeUpdatedAt,
				AnimeEpisodes: []*animeEpisode.AnimeEpisode{},
			}
		}

		// Add episode if it exists
		if result.EpisodeID != nil {
			episode := &animeEpisode.AnimeEpisode{
				ID:        *result.EpisodeID,
				AnimeID:   &result.AnimeID,
				Episode:   result.EpisodeNumber,
				TitleEn:   result.EpisodeTitleEn,
				TitleJp:   result.EpisodeTitleJp,
				Aired:     result.EpisodeAired,
				Synopsis:  result.EpisodeSynopsis,
				CreatedAt: *result.EpisodeCreatedAt,
				UpdatedAt: *result.EpisodeUpdatedAt,
			}
			animeMap[result.AnimeID].AnimeEpisodes = append(animeMap[result.AnimeID].AnimeEpisodes, episode)
		}
	}

	// Convert map to slice
	var animes []*Anime
	for _, anime := range animeMap {
		animes = append(animes, anime)
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})

	return animes, nil
}

// FindByIDs fetches multiple anime by their IDs in a single query
func (a *AnimeRepository) FindByIDs(ctx context.Context, ids []string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).Where("id IN ?", ids).Find(&animes).Error
	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

// FindByIDsWithEpisodes fetches multiple anime by their IDs with episodes preloaded in a single query
func (a *AnimeRepository) FindByIDsWithEpisodes(ctx context.Context, ids []string) ([]*Anime, error) {
	startTime := time.Now()

	var animes []*Anime

	// First, fetch all anime
	err := a.db.DB.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&animes).Error

	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Then, batch fetch all episodes for these anime in one query
	var episodes []animeEpisode.AnimeEpisode
	err = a.db.DB.WithContext(ctx).
		Where("anime_id IN ?", ids).
		Order("anime_id ASC, episode ASC").
		Find(&episodes).Error

	if err != nil {
		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_episodes",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Group episodes by anime_id
	episodeMap := make(map[string][]*animeEpisode.AnimeEpisode)
	for _, episode := range episodes {
		if episode.AnimeID != nil {
			ep := episode // Create a copy to avoid pointer issues
			episodeMap[*episode.AnimeID] = append(episodeMap[*episode.AnimeID], &ep)
		}
	}

	// Assign episodes to their respective anime
	for i, anime := range animes {
		if eps, exists := episodeMap[anime.ID]; exists {
			animes[i].AnimeEpisodes = eps
		}
	}

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})
	return animes, nil
}

// FindBySeasonWithEpisodesOptimized - Performance optimized version
func (a *AnimeRepository) FindBySeasonWithEpisodesOptimized(ctx context.Context, season string) ([]*Anime, error) {
	// Try cache first if available
	if a.cache != nil {
		key := fmt.Sprintf("%s:season_episodes_optimized:%s", a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1], season)
		var animeList []*Anime
		err := a.cache.GetJSON(ctx, key, &animeList)
		if err == nil {
			return animeList, nil
		}
		// Continue to database if cache miss or error
	}

	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeRepository.FindBySeasonWithEpisodesOptimized",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "repository"),
			attribute.String("anime.season", season),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).
		Table("anime_seasons as s").
		Joins("INNER JOIN anime as a ON s.anime_id = a.id").
		Joins("LEFT JOIN episodes as e ON a.id = e.anime_id").
		Where("s.season = ?", season).
		Find(&animes).Error

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			metrics.TableAnimeSeason,
			metrics.MethodSelect,
			metrics.Error,
		)
		return nil, err
	}

	duration := time.Since(startTime).Milliseconds()
	span.SetAttributes(
		attribute.Int64("duration_ms", duration),
		attribute.Int("result_count", len(animes)),
	)

	metrics.GetAppMetrics().DatabaseMetric(
		float64(duration),
		metrics.TableAnimeSeason,
		metrics.MethodSelect,
		metrics.Success,
	)

	// Store in cache if available
	if a.cache != nil {
		key := fmt.Sprintf("%s:season_episodes_optimized:%s", a.cache.GetKeyBuilder().AnimePattern()[:len(a.cache.GetKeyBuilder().AnimePattern())-1], season)
		_ = a.cache.SetJSON(ctx, key, animes, a.cache.GetSeasonTTL())
	}

	return animes, nil
}

// FindBySeasonAnimeOnlyOptimized - Ultra-fast version that only fetches anime without episodes
func (a *AnimeRepository) FindBySeasonAnimeOnlyOptimized(ctx context.Context, season string) ([]*Anime, error) {
	// Try cache first if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().AnimeBySeasonWithFields(season, nil)
		var animeList []*Anime
		err := a.cache.GetJSON(ctx, key, &animeList)
		if err == nil {
			return animeList, nil
		}
		// Continue to database if cache miss or error
	}

	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeRepository.FindBySeasonAnimeOnlyOptimized",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "repository"),
			attribute.String("anime.season", season),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	startTime := time.Now()

	var animes []*Anime
	err := a.db.DB.WithContext(ctx).
		Table("anime_seasons as s").
		Joins("INNER JOIN anime as a ON s.anime_id = a.id").
		Where("s.season = ?", season).
		Find(&animes).Error

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			metrics.TableAnimeSeason,
			metrics.MethodSelect,
			metrics.Error,
		)
		return nil, err
	}

	// Don't initialize episodes - let them be nil when not requested
	// This allows GraphQL to only include episodes field when explicitly queried

	duration := time.Since(startTime).Milliseconds()
	span.SetAttributes(
		attribute.Int64("duration_ms", duration),
		attribute.Int("result_count", len(animes)),
	)

	metrics.GetAppMetrics().DatabaseMetric(
		float64(duration),
		metrics.TableAnimeSeason,
		metrics.MethodSelect,
		metrics.Success,
	)

	// Store in cache if available
	if a.cache != nil {
		key := a.cache.GetKeyBuilder().AnimeBySeasonWithFields(season, nil)
		_ = a.cache.SetJSON(ctx, key, animes, a.cache.GetSeasonTTL())
	}

	return animes, nil
}

// FindBySeasonWithIndexHints - Uses MySQL index hints for better performance
func (a *AnimeRepository) FindBySeasonWithIndexHints(ctx context.Context, season string) ([]*Anime, error) {
	return a.FindBySeasonWithEpisodesOptimized(ctx, season)
}

// FindBySeasonBatched - Uses batched queries to avoid cartesian product
func (a *AnimeRepository) FindBySeasonBatched(ctx context.Context, season string) ([]*Anime, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeRepository.FindBySeasonBatched",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "repository"),
			attribute.String("anime.season", season),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	startTime := time.Now()

	// Step 1: Get anime IDs for the season
	var animeIDs []string
	err := a.db.DB.WithContext(ctx).
		Table("anime_seasons").
		Where("season = ?", season).
		Pluck("anime_id", &animeIDs).Error

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if len(animeIDs) == 0 {
		return []*Anime{}, nil
	}

	// Step 2: Get anime data
	var animes []*Anime
	err = a.db.DB.WithContext(ctx).
		Where("id IN ?", animeIDs).
		Find(&animes).Error

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Step 3: Get episodes for all anime in one query
	var episodes []*animeEpisode.AnimeEpisode
	err = a.db.DB.WithContext(ctx).
		Where("anime_id IN ?", animeIDs).
		Order("anime_id, episode").
		Find(&episodes).Error

	if err != nil {
		// Episodes are optional, continue without them
		span.SetAttributes(attribute.String("episodes_error", err.Error()))
	}

	// Step 4: Group episodes by anime ID
	episodeMap := make(map[string][]*animeEpisode.AnimeEpisode)
	for _, episode := range episodes {
		if episode.AnimeID != nil {
			episodeMap[*episode.AnimeID] = append(episodeMap[*episode.AnimeID], episode)
		}
	}

	// Step 5: Assign episodes to anime
	for _, a := range animes {
		if eps, exists := episodeMap[a.ID]; exists {
			a.AnimeEpisodes = eps
		} else {
			a.AnimeEpisodes = make([]*animeEpisode.AnimeEpisode, 0)
		}
	}

	duration := time.Since(startTime).Milliseconds()
	span.SetAttributes(
		attribute.Int64("duration_ms", duration),
		attribute.Int("result_count", len(animes)),
		attribute.Int("episode_count", len(episodes)),
	)

	metrics.GetAppMetrics().DatabaseMetric(
		float64(duration),
		metrics.TableAnimeSeason,
		metrics.MethodSelect,
		metrics.Success,
	)

	return animes, nil
}

// FindBySeasonWithFieldSelection - Optimized query that only selects requested fields
func (a *AnimeRepository) FindBySeasonWithFieldSelection(ctx context.Context, season string, fieldSelection *FieldSelection, limit int) ([]*Anime, error) {
	tracer := tracing.GetTracer(ctx)

	// Try cache first if available
	if a.cache != nil {
		ctx, cacheSpan := tracer.Start(ctx, "AnimeRepository.CacheGet",
			trace.WithAttributes(
				attribute.String("cache.operation", "get"),
				attribute.String("cache.layer", "repository"),
				attribute.String("anime.season", season),
				attribute.Int("cache.field_count", len(fieldSelection.Fields)),
			),
			trace.WithSpanKind(trace.SpanKindInternal),
			tracing.GetEnvironmentAttribute(),
		)

		// Convert map[string]bool to []string for cache key
		var fields []string
		for field := range fieldSelection.Fields {
			fields = append(fields, field)
		}
		key := a.cache.GetKeyBuilder().AnimeBySeasonWithFields(season, fields)
		cacheSpan.SetAttributes(
			attribute.String("cache.key", key),
			attribute.StringSlice("cache.fields", fields),
		)

		cacheStartTime := time.Now()
		var animeList []*Anime
		err := a.cache.GetJSON(ctx, key, &animeList)
		cacheDuration := time.Since(cacheStartTime)

		cacheSpan.SetAttributes(
			attribute.Int64("cache.duration_us", cacheDuration.Microseconds()),
			attribute.Int64("cache.duration_ms", cacheDuration.Milliseconds()),
		)

		if err == nil {
			// Apply limit to cached results if specified
			if limit > 0 && len(animeList) > limit {
				animeList = animeList[:limit]
			}
			cacheSpan.SetAttributes(
				attribute.String("cache.result", "hit"),
				attribute.Int("cache.items_returned", len(animeList)),
			)
			cacheSpan.SetStatus(codes.Ok, "cache hit")
			cacheSpan.End()
			return animeList, nil
		} else {
			cacheSpan.SetAttributes(attribute.String("cache.result", "miss"))
			if err != cache.ErrCacheMiss {
				cacheSpan.RecordError(err)
				cacheSpan.SetStatus(codes.Error, err.Error())
			} else {
				cacheSpan.SetStatus(codes.Ok, "cache miss")
			}
			cacheSpan.End()
		}
		// Continue to database if cache miss or error
	}

	ctx, span := tracer.Start(ctx, "AnimeRepository.FindBySeasonWithFieldSelection",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "repository"),
			attribute.String("anime.season", season),
			attribute.Int("anime.fields_count", len(fieldSelection.Fields)),
			attribute.Int("anime.limit", limit),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	startTime := time.Now()

	// Build the select clause based on requested fields
	selectClause := fieldSelection.BuildSelectClause("anime")

	// Use raw SQL for maximum performance and field selection control
	var animeList []*Anime
	query := `
		SELECT ` + selectClause + `
		FROM anime
		INNER JOIN anime_seasons ON anime.id = anime_seasons.anime_id
		WHERE anime_seasons.season = ?
		ORDER BY anime.ranking ASC, anime.title_en ASC
		LIMIT ?`

	span.SetAttributes(attribute.String("db.query", query))
	result := a.db.DB.WithContext(ctx).Raw(query, season, limit).Scan(&animeList)
	if result.Error != nil {
		metrics.GetAppMetrics().DatabaseMetric(
			float64(time.Since(startTime).Milliseconds()),
			metrics.TableAnimeSeason,
			metrics.MethodSelect,
			metrics.Error,
		)
		return nil, result.Error
	}

	span.SetAttributes(
		attribute.Int("selected_fields_count", len(fieldSelection.Fields)),
		attribute.Int("result_count", len(animeList)),
	)

	metrics.GetAppMetrics().DatabaseMetric(
		float64(time.Since(startTime).Milliseconds()),
		metrics.TableAnimeSeason,
		metrics.MethodSelect,
		metrics.Success,
	)

	// Store in cache if available
	if a.cache != nil {
		// Convert map[string]bool to []string for cache key
		var fields []string
		for field := range fieldSelection.Fields {
			fields = append(fields, field)
		}
		key := a.cache.GetKeyBuilder().AnimeBySeasonWithFields(season, fields)
		_ = a.cache.SetJSON(ctx, key, animeList, a.cache.GetSeasonTTL())
	}

	return animeList, nil
}

// Start-of-day in a specific zone (keeps date math stable)
func startOfDayIn(t time.Time, loc *time.Location) time.Time {
	y, m, d := t.In(loc).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, loc)
}
