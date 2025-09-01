package anime_season

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_season"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type AnimeSeasonServiceImpl interface {
	FindByAnimeID(ctx context.Context, animeID string) ([]*anime_season.AnimeSeason, error)
	FindBySeason(ctx context.Context, season string) ([]*anime_season.AnimeSeason, error)
	Create(ctx context.Context, animeSeason *anime_season.AnimeSeason) error
	Update(ctx context.Context, animeSeason *anime_season.AnimeSeason) error
	Delete(ctx context.Context, id string) error
}

type AnimeSeasonService struct {
	Repository anime_season.AnimeSeasonRepositoryImpl
}

func NewAnimeSeasonService(repository anime_season.AnimeSeasonRepositoryImpl) AnimeSeasonServiceImpl {
	return &AnimeSeasonService{
		Repository: repository,
	}
}

func (s *AnimeSeasonService) FindByAnimeID(ctx context.Context, animeID string) ([]*anime_season.AnimeSeason, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "FindByAnimeID")
	span.SetTag("service", "anime_season")
	span.SetTag("type", "service")
	defer span.Finish()

	return s.Repository.FindByAnimeID(spanCtx, animeID)
}

func (s *AnimeSeasonService) FindBySeason(ctx context.Context, season string) ([]*anime_season.AnimeSeason, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "FindBySeason")
	span.SetTag("service", "anime_season")
	span.SetTag("type", "service")
	defer span.Finish()

	return s.Repository.FindBySeason(spanCtx, season)
}

func (s *AnimeSeasonService) Create(ctx context.Context, animeSeason *anime_season.AnimeSeason) error {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "CreateAnimeSeason")
	span.SetTag("service", "anime_season")
	span.SetTag("type", "service")
	defer span.Finish()

	return s.Repository.Create(spanCtx, animeSeason)
}

func (s *AnimeSeasonService) Update(ctx context.Context, animeSeason *anime_season.AnimeSeason) error {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "UpdateAnimeSeason")
	span.SetTag("service", "anime_season")
	span.SetTag("type", "service")
	defer span.Finish()

	return s.Repository.Update(spanCtx, animeSeason)
}

func (s *AnimeSeasonService) Delete(ctx context.Context, id string) error {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "DeleteAnimeSeason")
	span.SetTag("service", "anime_season")
	span.SetTag("type", "service")
	defer span.Finish()

	return s.Repository.Delete(spanCtx, id)
}