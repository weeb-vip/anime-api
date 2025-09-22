package resolvers

import (
	"context"
	"testing"
	"time"

	"github.com/weeb-vip/anime-api/graph/model"
	anime_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	anime_episode_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	anime_service "github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime_season"
	"go.uber.org/mock/gomock"
)

func BenchmarkAnimeBySeasonsOriginal(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockAnimeService := NewMockAnimeService(ctrl)
	mockAnimeSeasonService := NewMockAnimeSeasonService(ctrl)

	ctx := context.Background()
	season := "SPRING_2024"

	// Create test data
	now := time.Now()
	testAnime := &anime_repo.Anime{
		ID:       "anime-1",
		TitleEn:  stringPtr("Test Anime"),
		Episodes: intPtr(12),
		AnimeEpisodes: []*anime_episode_repo.AnimeEpisode{
			{
				ID:        "ep-1",
				AnimeID:   stringPtr("anime-1"),
				Episode:   intPtr(1),
				TitleEn:   stringPtr("Episode 1"),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Mock the original method
	mockAnimeService.EXPECT().
		AnimeBySeasonWithEpisodes(gomock.Any(), season).
		Return([]*anime_repo.Anime{testAnime}, nil).
		AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use the old method for comparison
		_, err := animeBySeasonsOriginal(ctx, mockAnimeSeasonService, mockAnimeService, season)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAnimeBySeasonsOptimized(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockAnimeService := NewMockAnimeService(ctrl)
	mockAnimeSeasonService := NewMockAnimeSeasonService(ctrl)

	ctx := context.Background()
	season := "SPRING_2024"

	// Create test data
	now := time.Now()
	testAnime := &anime_repo.Anime{
		ID:       "anime-1",
		TitleEn:  stringPtr("Test Anime"),
		Episodes: intPtr(12),
		AnimeEpisodes: []*anime_episode_repo.AnimeEpisode{
			{
				ID:        "ep-1",
				AnimeID:   stringPtr("anime-1"),
				Episode:   intPtr(1),
				TitleEn:   stringPtr("Episode 1"),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Mock the optimized method
	mockAnimeService.EXPECT().
		AnimeBySeasonWithEpisodesOptimized(gomock.Any(), season).
		Return([]*anime_repo.Anime{testAnime}, nil).
		AnyTimes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := AnimeBySeasons(ctx, mockAnimeSeasonService, mockAnimeService, season)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Original method for comparison
func animeBySeasonsOriginal(ctx context.Context, animeSeasonService anime_season.AnimeSeasonServiceImpl, animeService anime_service.AnimeServiceImpl, season string) ([]*model.Anime, error) {
	animeList, err := animeService.AnimeBySeasonWithEpisodes(ctx, season)
	if err != nil {
		return nil, err
	}

	var result []*model.Anime
	for _, animeEntity := range animeList {
		animeGraphQL, err := transformAnimeToGraphQL(*animeEntity)
		if err != nil {
			continue
		}
		result = append(result, animeGraphQL)
	}

	return result, nil
}

func TestOptimizedQueryPerformance(t *testing.T) {
	// This test is mainly to document the performance improvements
	// In real usage, you should see significant improvements:
	// 1. Reduced SQL query complexity
	// 2. Fewer fields transferred
	// 3. Better memory allocation patterns
	// 4. Reduced network overhead between app and database

	t.Log("Performance optimizations implemented:")
	t.Log("1. Selective field loading (only required fields)")
	t.Log("2. Removed expensive ORDER BY on episode numbers")
	t.Log("3. Pre-allocated slices and maps for better memory usage")
	t.Log("4. Added proper tracing for performance monitoring")
	t.Log("5. Option for anime-only queries when episodes aren't needed")
}