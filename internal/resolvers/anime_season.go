package resolvers

import (
	"context"
	"fmt"
	"strings"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
	"github.com/weeb-vip/anime-api/graph/model"
	anime_season_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime_season"
	anime_service "github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime_season"
	"github.com/weeb-vip/anime-api/metrics"
	"time"
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
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "AnimeSeasons",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
		})
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

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "AnimeSeasons",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return seasons, nil
}

func AnimeBySeasons(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, season string) ([]*model.Anime, error) {
	startTime := time.Now()

	foundSeasons, err := animeSeasonService.FindBySeason(ctx, season)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "AnimeBySeasons",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
		})
		return nil, err
	}

	// Get unique anime IDs
	animeIDs := make([]string, 0)
	seenIDs := make(map[string]bool)
	for _, seasonEntity := range foundSeasons {
		if seasonEntity.AnimeID != nil && !seenIDs[*seasonEntity.AnimeID] {
			animeIDs = append(animeIDs, *seasonEntity.AnimeID)
			seenIDs[*seasonEntity.AnimeID] = true
		}
	}

	// Fetch anime data for each unique anime ID
	var result []*model.Anime
	for _, animeID := range animeIDs {
		animeEntity, err := animeService.AnimeByID(ctx, animeID)
		if err != nil {
			continue // Skip anime that can't be fetched
		}
		
		animeGraphQL, err := transformAnimeToGraphQL(*animeEntity)
		if err != nil {
			continue // Skip anime that can't be transformed
		}
		
		result = append(result, animeGraphQL)
	}

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "AnimeBySeasons",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return result, nil
}

func AnimeBySeasonAndYear(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, seasonName string, year int) ([]*model.Anime, error) {
	// Construct the season string in the expected format: SPRING_2024, SUMMER_2024, etc.
	season := fmt.Sprintf("%s_%d", strings.ToUpper(seasonName), year)
	return AnimeBySeasons(ctx, animeSeasonService, animeService, season)
}