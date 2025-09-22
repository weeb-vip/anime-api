package anime

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/tracing"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"time"
)

type AnimeServiceImpl interface {
	AnimeByID(ctx context.Context, id string) (*anime.Anime, error)
	AnimeByIDWithEpisodes(ctx context.Context, id string) (*anime.Anime, error)
	AnimeByIDs(ctx context.Context, ids []string) ([]*anime.Anime, error)
	AnimeByIDsWithEpisodes(ctx context.Context, ids []string) ([]*anime.Anime, error)
	TopRatedAnime(ctx context.Context, limit int) ([]*anime.Anime, error)
	TopRatedAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error)
	MostPopularAnime(ctx context.Context, limit int) ([]*anime.Anime, error)
	MostPopularAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error)
	NewestAnime(ctx context.Context, limit int) ([]*anime.Anime, error)
	NewestAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error)
	AiringAnime(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime.AnimeWithNextEpisode, error)
	AiringAnimeWithEpisodes(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime.Anime, error)
	SearchedAnime(ctx context.Context, query string, page int, limit int) ([]*anime.Anime, error)
	SearchedAnimeWithEpisodes(ctx context.Context, query string, page int, limit int) ([]*anime.Anime, error)
	AnimeBySeasonWithEpisodes(ctx context.Context, season string) ([]*anime.Anime, error)
	AnimeBySeasonWithEpisodesOptimized(ctx context.Context, season string) ([]*anime.Anime, error)
	AnimeBySeasonWithIndexHints(ctx context.Context, season string) ([]*anime.Anime, error)
	AnimeBySeasonBatched(ctx context.Context, season string) ([]*anime.Anime, error)
	AnimeBySeasonOptimized(ctx context.Context, season string) ([]*anime.Anime, error)
}

type AnimeService struct {
	Repository anime.AnimeRepositoryImpl
}

func NewAnimeService(animeRepository anime.AnimeRepositoryImpl) AnimeServiceImpl {
	return &AnimeService{
		Repository: animeRepository,
	}
}

func (a *AnimeService) AnimeByID(ctx context.Context, id string) (*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeByID")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.FindById(spanCtx, id)
}

func (a *AnimeService) TopRatedAnime(ctx context.Context, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "TopRatedAnime")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.TopRatedAnime(spanCtx, limit)
}

func (a *AnimeService) MostPopularAnime(ctx context.Context, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "MostPopularAnime")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.MostPopularAnime(spanCtx, limit)
}

func (a *AnimeService) NewestAnime(ctx context.Context, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "NewestAnime")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.NewestAnime(spanCtx, limit)
}

func (a *AnimeService) AiringAnime(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime.AnimeWithNextEpisode, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AiringAnime")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	if startDate == nil && endDate == nil && days == nil {
		return a.Repository.AiringAnime(spanCtx)
	}
	if endDate == nil {
		return a.Repository.AiringAnimeDays(spanCtx, startDate, days)
	}
	return a.Repository.AiringAnimeEndDate(spanCtx, startDate, endDate)
}

func (a *AnimeService) SearchedAnime(ctx context.Context, query string, page int, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "SearchedAnime")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.SearchAnime(spanCtx, query, page, limit)
}

func (a *AnimeService) AnimeByIDWithEpisodes(ctx context.Context, id string) (*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeByIDWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.FindByIdWithEpisodes(spanCtx, id)
}

func (a *AnimeService) TopRatedAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "TopRatedAnimeWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.TopRatedAnimeWithEpisodes(spanCtx, limit)
}

func (a *AnimeService) MostPopularAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "MostPopularAnimeWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.MostPopularAnimeWithEpisodes(spanCtx, limit)
}

func (a *AnimeService) NewestAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "NewestAnimeWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.NewestAnimeWithEpisodes(spanCtx, limit)
}

func (a *AnimeService) SearchedAnimeWithEpisodes(ctx context.Context, query string, page int, limit int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "SearchedAnimeWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.SearchAnimeWithEpisodes(spanCtx, query, page, limit)
}

func (a *AnimeService) AiringAnimeWithEpisodes(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AiringAnimeWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.AiringAnimeWithEpisodes(spanCtx, startDate, endDate, days)
}

func (a *AnimeService) AnimeBySeasonWithEpisodes(ctx context.Context, season string) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeBySeasonWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.FindBySeasonWithEpisodes(spanCtx, season)
}

func (a *AnimeService) AnimeBySeasonWithEpisodesOptimized(ctx context.Context, season string) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeBySeasonWithEpisodesOptimized")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.FindBySeasonWithEpisodesOptimized(spanCtx, season)
}


func (a *AnimeService) AnimeBySeasonWithIndexHints(ctx context.Context, season string) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeBySeasonWithIndexHints")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.FindBySeasonWithIndexHints(spanCtx, season)
}

func (a *AnimeService) AnimeBySeasonBatched(ctx context.Context, season string) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeBySeasonBatched")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.FindBySeasonBatched(spanCtx, season)
}

func (a *AnimeService) AnimeBySeasonOptimized(ctx context.Context, season string) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeBySeasonOptimized")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("environment", tracing.GetEnvironmentTag())
	defer span.Finish()

	return a.Repository.FindBySeasonAnimeOnlyOptimized(spanCtx, season)
}
