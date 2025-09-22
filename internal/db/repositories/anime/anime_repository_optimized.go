package anime

import (
	"context"
	"time"

	"github.com/weeb-vip/anime-api/metrics"
	"github.com/weeb-vip/anime-api/tracing"
	animeEpisode "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// FindBySeasonWithEpisodesOptimized - Performance optimized version that:
// 1. Only selects necessary fields (not all columns)
// 2. Uses indexes efficiently
// 3. Minimizes data transfer
// 4. Optimizes the JOIN strategy
func (a *AnimeRepository) FindBySeasonWithEpisodesOptimized(ctx context.Context, season string) ([]*Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "FindBySeasonWithEpisodesOptimized")
	span.SetTag("service", "anime")
	span.SetTag("type", "repository")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	span.SetTag("season", season)
	defer span.Finish()

	startTime := time.Now()

	// Optimized query structure to reduce fields and improve performance
	type OptimizedAnimeResult struct {
		// Core anime fields only (not all fields)
		AnimeID         string    `gorm:"column:anime_id"`
		TitleEn         *string   `gorm:"column:title_en"`
		TitleJp         *string   `gorm:"column:title_jp"`
		TitleRomaji     *string   `gorm:"column:title_romaji"`
		TitleKanji      *string   `gorm:"column:title_kanji"`
		TitleSynonyms   *string   `gorm:"column:title_synonyms"`
		ImageURL        *string   `gorm:"column:image_url"`
		Synopsis        *string   `gorm:"column:synopsis"`
		Episodes        *int      `gorm:"column:episodes"`
		Status          *string   `gorm:"column:status"`
		StartDate       *string   `gorm:"column:start_date"`
		EndDate         *string   `gorm:"column:end_date"`
		Genres          *string   `gorm:"column:genres"`
		Duration        *string   `gorm:"column:duration"`
		Broadcast       *string   `gorm:"column:broadcast"`
		Source          *string   `gorm:"column:source"`
		Licensors       *string   `gorm:"column:licensors"`
		Studios         *string   `gorm:"column:studios"`
		Rating          *string   `gorm:"column:rating"`
		Ranking         *int      `gorm:"column:ranking"`
		AnimeCreatedAt  time.Time `gorm:"column:anime_created_at"`
		AnimeUpdatedAt  time.Time `gorm:"column:anime_updated_at"`

		// Essential episode fields only
		EpisodeID       *string    `gorm:"column:episode_id"`
		EpisodeNumber   *int       `gorm:"column:episode"`
		EpisodeTitleEn  *string    `gorm:"column:episode_title_en"`
		EpisodeTitleJp  *string    `gorm:"column:episode_title_jp"`
		EpisodeAired    *time.Time `gorm:"column:aired"`
		EpisodeSynopsis *string    `gorm:"column:episode_synopsis"`
		EpisodeCreatedAt *time.Time `gorm:"column:episode_created_at"`
		EpisodeUpdatedAt *time.Time `gorm:"column:episode_updated_at"`
	}

	var results []OptimizedAnimeResult

	// Strategy 1: Optimized field selection with proper indexing hints
	err := a.db.DB.WithContext(spanCtx).
		Select(`
			a.id as anime_id,
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
		// Remove expensive ordering for now - we can sort in Go if needed
		Find(&results).Error

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.msg", err.Error())

		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Optimized grouping logic
	animeMap := make(map[string]*Anime, len(results)/10) // Pre-allocate with estimated size
	for _, result := range results {
		// Create anime if not exists
		if _, exists := animeMap[result.AnimeID]; !exists {
			animeMap[result.AnimeID] = &Anime{
				ID:            result.AnimeID,
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
				AnimeEpisodes: make([]*animeEpisode.AnimeEpisode, 0), // Pre-allocate empty slice
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

	// Convert map to slice with pre-allocation
	animes := make([]*Anime, 0, len(animeMap))
	for _, anime := range animeMap {
		animes = append(animes, anime)
	}

	duration := time.Since(startTime).Milliseconds()
	span.SetTag("duration_ms", duration)
	span.SetTag("result_count", len(animes))

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(duration), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})

	return animes, nil
}

// FindBySeasonAnimeOnlyOptimized - Ultra-fast version that only fetches anime without episodes
// Use this if episodes aren't needed for the query
func (a *AnimeRepository) FindBySeasonAnimeOnlyOptimized(ctx context.Context, season string) ([]*Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "FindBySeasonAnimeOnlyOptimized")
	span.SetTag("service", "anime")
	span.SetTag("type", "repository")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	span.SetTag("season", season)
	defer span.Finish()

	startTime := time.Now()

	var animes []*Anime

	// Much faster query - no episodes join
	err := a.db.DB.WithContext(spanCtx).
		Select(`
			a.id,
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
			a.created_at,
			a.updated_at
		`).
		Table("anime_seasons as s").
		Joins("INNER JOIN anime as a ON s.anime_id = a.id").
		Where("s.season = ?", season).
		Find(&animes).Error

	if err != nil {
		span.SetTag("error", true)
		span.SetTag("error.msg", err.Error())

		_ = metrics.NewMetricsInstance().DatabaseMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.DatabaseMetricLabels{
			Service: "anime-api",
			Table:   "anime_seasons",
			Method:  metrics_lib.DatabaseMetricMethodSelect,
			Result:  metrics_lib.Error,
			Env:     metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Initialize empty episodes for each anime to prevent nil pointer issues
	for _, anime := range animes {
		anime.AnimeEpisodes = make([]*animeEpisode.AnimeEpisode, 0)
	}

	duration := time.Since(startTime).Milliseconds()
	span.SetTag("duration_ms", duration)
	span.SetTag("result_count", len(animes))

	_ = metrics.NewMetricsInstance().DatabaseMetric(float64(duration), metrics_lib.DatabaseMetricLabels{
		Service: "anime-api",
		Table:   "anime_seasons",
		Method:  metrics_lib.DatabaseMetricMethodSelect,
		Result:  metrics_lib.Success,
		Env:     metrics.GetCurrentEnv(),
	})

	return animes, nil
}