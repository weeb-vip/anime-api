package resolvers

import (
	"context"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
	"github.com/weeb-vip/anime-api/graph/model"
	anime_service "github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime_season"
	"github.com/weeb-vip/anime-api/metrics"
	"time"
)

// AnimeBySeasonsOptimized uses batch fetching instead of N+1 queries
func AnimeBySeasonsOptimized(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, season string) ([]*model.Anime, error) {
	startTime := time.Now()

	foundSeasons, err := animeSeasonService.FindBySeason(ctx, season)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "AnimeBySeasonsOptimized",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
			Env:      metrics.GetCurrentEnv(),
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

	// Batch fetch all anime at once using FindByIDs (needs to be added to service/repository)
	// For now, we'll add a new method to fetch multiple anime by IDs
	animeList, err := animeService.AnimeByIDsWithEpisodes(ctx, animeIDs)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "AnimeBySeasonsOptimized",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
			Env:      metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Transform to GraphQL models
	var result []*model.Anime
	for _, animeEntity := range animeList {
		animeGraphQL, err := transformAnimeToGraphQL(*animeEntity)
		if err != nil {
			continue // Skip anime that can't be transformed
		}
		result = append(result, animeGraphQL)
	}

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "AnimeBySeasonsOptimized",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
		Env:      metrics.GetCurrentEnv(),
	})

	return result, nil
}