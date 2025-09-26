package anime

import (
	"context"
	"time"

	"github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	AnimeBySeasonWithFieldSelection(ctx context.Context, season string, fields *anime.FieldSelection) ([]*anime.Anime, error)
}

type AnimeService struct {
	Repository anime.AnimeRepositoryImpl
}

// startServiceSpan starts a new OpenTelemetry span for service operations
func (a *AnimeService) startServiceSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
	tracer := tracing.GetTracer(ctx)
	return tracer.Start(ctx, "AnimeService."+operationName,
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "service"),
			attribute.String("method", operationName),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
		tracing.GetEnvironmentAttribute(),
	)
}

func NewAnimeService(animeRepository anime.AnimeRepositoryImpl) AnimeServiceImpl {
	return &AnimeService{
		Repository: animeRepository,
	}
}

func (a *AnimeService) AnimeByID(ctx context.Context, id string) (*anime.Anime, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeService.AnimeByID",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "service"),
			attribute.String("anime.id", id),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	return a.Repository.FindById(ctx, id)
}

func (a *AnimeService) TopRatedAnime(ctx context.Context, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "TopRatedAnime")
	defer span.End()

	return a.Repository.TopRatedAnime(ctx, limit)
}

func (a *AnimeService) MostPopularAnime(ctx context.Context, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "MostPopularAnime")
	defer span.End()

	return a.Repository.MostPopularAnime(ctx, limit)
}

func (a *AnimeService) NewestAnime(ctx context.Context, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "NewestAnime")
	defer span.End()

	return a.Repository.NewestAnime(ctx, limit)
}

func (a *AnimeService) AiringAnime(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime.AnimeWithNextEpisode, error) {
	ctx, span := a.startServiceSpan(ctx, "AiringAnime")
	defer span.End()

	if startDate == nil && endDate == nil && days == nil {
		return a.Repository.AiringAnime(ctx)
	}
	if endDate == nil {
		return a.Repository.AiringAnimeDays(ctx, startDate, days)
	}
	return a.Repository.AiringAnimeEndDate(ctx, startDate, endDate)
}

func (a *AnimeService) SearchedAnime(ctx context.Context, query string, page int, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "SearchedAnime")
	span.SetAttributes(
		attribute.String("search.query", query),
		attribute.Int("pagination.page", page),
		attribute.Int("pagination.limit", limit),
	)
	defer span.End()

	return a.Repository.SearchAnime(ctx, query, page, limit)
}

func (a *AnimeService) AnimeByIDWithEpisodes(ctx context.Context, id string) (*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "AnimeByIDWithEpisodes")
	span.SetAttributes(attribute.String("anime.id", id))
	defer span.End()

	return a.Repository.FindByIdWithEpisodes(ctx, id)
}

func (a *AnimeService) TopRatedAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "TopRatedAnimeWithEpisodes")
	defer span.End()

	return a.Repository.TopRatedAnimeWithEpisodes(ctx, limit)
}

func (a *AnimeService) MostPopularAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "MostPopularAnimeWithEpisodes")
	defer span.End()

	return a.Repository.MostPopularAnimeWithEpisodes(ctx, limit)
}

func (a *AnimeService) NewestAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "NewestAnimeWithEpisodes")
	defer span.End()

	return a.Repository.NewestAnimeWithEpisodes(ctx, limit)
}

func (a *AnimeService) SearchedAnimeWithEpisodes(ctx context.Context, query string, page int, limit int) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "SearchedAnimeWithEpisodes")
	span.SetAttributes(
		attribute.String("search.query", query),
		attribute.Int("pagination.page", page),
		attribute.Int("pagination.limit", limit),
	)
	defer span.End()

	return a.Repository.SearchAnimeWithEpisodes(ctx, query, page, limit)
}

func (a *AnimeService) AiringAnimeWithEpisodes(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime.Anime, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeService.AiringAnimeWithEpisodes",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "service"),
			attribute.String("method", "AiringAnimeWithEpisodes"),
		),
		trace.WithSpanKind(trace.SpanKindInternal),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	return a.Repository.AiringAnimeWithEpisodes(ctx, startDate, endDate, days)
}

func (a *AnimeService) AnimeBySeasonWithEpisodes(ctx context.Context, season string) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "AnimeBySeasonWithEpisodes")
	span.SetAttributes(attribute.String("anime.season", season))
	defer span.End()

	return a.Repository.FindBySeasonWithEpisodes(ctx, season)
}

func (a *AnimeService) AnimeBySeasonWithEpisodesOptimized(ctx context.Context, season string) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "AnimeBySeasonWithEpisodesOptimized")
	span.SetAttributes(attribute.String("anime.season", season))
	defer span.End()

	return a.Repository.FindBySeasonWithEpisodesOptimized(ctx, season)
}


func (a *AnimeService) AnimeBySeasonWithIndexHints(ctx context.Context, season string) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "AnimeBySeasonWithIndexHints")
	span.SetAttributes(attribute.String("anime.season", season))
	defer span.End()

	return a.Repository.FindBySeasonWithIndexHints(ctx, season)
}

func (a *AnimeService) AnimeBySeasonBatched(ctx context.Context, season string) ([]*anime.Anime, error) {
	ctx, span := a.startServiceSpan(ctx, "AnimeBySeasonBatched")
	span.SetAttributes(attribute.String("anime.season", season))
	defer span.End()

	return a.Repository.FindBySeasonBatched(ctx, season)
}

func (a *AnimeService) AnimeBySeasonOptimized(ctx context.Context, season string) ([]*anime.Anime, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeService.AnimeBySeasonOptimized",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "service"),
			attribute.String("anime.season", season),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	return a.Repository.FindBySeasonAnimeOnlyOptimized(ctx, season)
}

func (a *AnimeService) AnimeBySeasonWithFieldSelection(ctx context.Context, season string, fields *anime.FieldSelection) ([]*anime.Anime, error) {
	tracer := tracing.GetTracer(ctx)
	ctx, span := tracer.Start(ctx, "AnimeService.AnimeBySeasonWithFieldSelection",
		trace.WithAttributes(
			attribute.String("service", "anime"),
			attribute.String("type", "service"),
			attribute.String("anime.season", season),
			attribute.Int("anime.fields_count", len(fields.Fields)),
		),
		tracing.GetEnvironmentAttribute(),
	)
	defer span.End()

	return a.Repository.FindBySeasonWithFieldSelection(ctx, season, fields)
}
