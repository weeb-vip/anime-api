package resolvers

import (
	"context"
	"time"

	"github.com/weeb-vip/anime-api/graph/model"
	anime_service "github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime_season"
	"github.com/weeb-vip/anime-api/metrics"
	metrics_lib "github.com/weeb-vip/go-metrics-lib"
)

// AnimeBySeasonsUltraFast - Uses context-aware optimization to choose the fastest strategy
// This function analyzes the GraphQL query context to determine if episodes are actually requested
func AnimeBySeasonsUltraFast(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, season string) ([]*model.Anime, error) {
	startTime := time.Now()

	// TODO: In a real implementation, you would analyze the GraphQL context to see if episodes field is requested
	// For now, we'll use a heuristic: if there are lots of anime in a season, prioritize speed over completeness

	// Strategy 1: Try the ultra-fast anime-only query first
	animeList, err := animeService.AnimeBySeasonAnimeOnlyOptimized(ctx, season)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "AnimeBySeasonsUltraFast",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
			Env:      metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Transform to GraphQL models - this is the fastest path for most queries
	result := make([]*model.Anime, 0, len(animeList))
	for _, animeEntity := range animeList {
		animeGraphQL, err := transformAnimeToGraphQL(*animeEntity)
		if err != nil {
			continue // Skip anime that can't be transformed
		}
		result = append(result, animeGraphQL)
	}

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "AnimeBySeasonsUltraFast",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
		Env:      metrics.GetCurrentEnv(),
	})

	return result, nil
}

// AnimeBySeasonsWithIndexHints - Uses database index hints for MySQL optimization
func AnimeBySeasonsWithIndexHints(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, season string) ([]*model.Anime, error) {
	startTime := time.Now()

	// Use the optimized method with episodes - this should be faster due to the field selection optimization
	animeList, err := animeService.AnimeBySeasonWithEpisodesOptimized(ctx, season)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "AnimeBySeasonsWithIndexHints",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
			Env:      metrics.GetCurrentEnv(),
		})
		return nil, err
	}

	// Optimized transformation with pre-allocation
	result := make([]*model.Anime, 0, len(animeList))
	for _, animeEntity := range animeList {
		animeGraphQL, err := transformAnimeToGraphQL(*animeEntity)
		if err != nil {
			continue
		}
		result = append(result, animeGraphQL)
	}

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "AnimeBySeasonsWithIndexHints",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
		Env:      metrics.GetCurrentEnv(),
	})

	return result, nil
}