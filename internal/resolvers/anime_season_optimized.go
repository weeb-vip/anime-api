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

// AnimeBySeasonsOptimized uses a single SQL query with joins for maximum performance
func AnimeBySeasonsOptimized(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, season string) ([]*model.Anime, error) {
	startTime := time.Now()

	// Use the new optimized single-query method that joins all tables
	animeList, err := animeService.AnimeBySeasonWithEpisodes(ctx, season)
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