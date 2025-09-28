package resolvers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/weeb-vip/anime-api/graph/model"
	anime_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	anime_season_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime_season"
	anime_service "github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime_season"
	"github.com/weeb-vip/anime-api/metrics"
)

func transformAnimeSeasonToGraphQL(animeSeasonEntity anime_season_repo.AnimeSeason) (*model.AnimeSeason, error) {
	season := animeSeasonEntity.Season // Now Season is a string scalar
	status := model.AnimeSeasonStatus(animeSeasonEntity.Status)

	return &model.AnimeSeason{
		ID:           animeSeasonEntity.ID,
		Season:       season,
		Status:       status,
		EpisodeCount: animeSeasonEntity.EpisodeCount,
		Notes:        animeSeasonEntity.Notes,
		AnimeID:      animeSeasonEntity.AnimeID,
		CreatedAt:    animeSeasonEntity.CreatedAt,
		UpdatedAt:    animeSeasonEntity.UpdatedAt,
	}, nil
}

func AnimeSeasons(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeID string) ([]*model.AnimeSeason, error) {
	startTime := time.Now()

	foundSeasons, err := animeSeasonService.FindByAnimeID(ctx, animeID)
	if err != nil {
		metrics.GetAppMetrics().ResolverMetric(
			float64(time.Since(startTime).Milliseconds()),
			"AnimeSeasons",
			metrics.Error,
		)
		return nil, err
	}

	var seasons []*model.AnimeSeason
	for _, seasonEntity := range foundSeasons {
		season, err := transformAnimeSeasonToGraphQL(*seasonEntity)
		if err != nil {
			return nil, err
		}
		seasons = append(seasons, season)
	}

	metrics.GetAppMetrics().ResolverMetric(
		float64(time.Since(startTime).Milliseconds()),
		"AnimeSeasons",
		metrics.Success,
	)

	return seasons, nil
}

func AnimeBySeasons(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, season string, limit *int) ([]*model.Anime, error) {
	startTime := time.Now()

	// Set default limit to 10 if not provided
	queryLimit := 10
	if limit != nil {
		queryLimit = *limit
	}

	// Extract field selection from GraphQL context to optimize query
	fieldSelection := ExtractAnimeFieldSelection(ctx)

	var animeList []*anime_repo.Anime
	var err error

	if fieldSelection != nil {
		// Use field-optimized query that only selects requested fields
		animeList, err = animeService.AnimeBySeasonWithFieldSelection(ctx, season, fieldSelection, queryLimit)
	} else {
		// Fall back to standard optimized method (no episodes)
		animeList, err = animeService.AnimeBySeasonOptimized(ctx, season)
		// Apply limit manually since this method doesn't support it
		if err == nil && len(animeList) > queryLimit {
			animeList = animeList[:queryLimit]
		}
	}

	if err != nil {
		metrics.GetAppMetrics().ResolverMetric(
			float64(time.Since(startTime).Milliseconds()),
			"AnimeBySeasons",
			metrics.Error,
		)
		return nil, err
	}

	// Transform to GraphQL models - pre-allocate for performance
	result := make([]*model.Anime, 0, len(animeList))
	for _, animeEntity := range animeList {
		animeGraphQL, err := transformAnimeToGraphQL(*animeEntity)
		if err != nil {
			continue // Skip anime that can't be transformed
		}
		result = append(result, animeGraphQL)
	}

	metrics.GetAppMetrics().ResolverMetric(
		float64(time.Since(startTime).Milliseconds()),
		"AnimeBySeasons",
		metrics.Success,
	)

	return result, nil
}

func AnimeBySeasonAndYear(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, seasonName string, year int, limit *int) ([]*model.Anime, error) {
	// Construct the season string in the expected format: SPRING_2024, SUMMER_2024, etc.
	season := fmt.Sprintf("%s_%d", strings.ToUpper(seasonName), year)
	return AnimeBySeasons(ctx, animeSeasonService, animeService, season, limit)
}