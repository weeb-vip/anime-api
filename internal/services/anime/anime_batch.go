package anime

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// AnimeByIDsWithEpisodes fetches multiple anime by their IDs in a single query with episodes preloaded
func (a *AnimeService) AnimeByIDsWithEpisodes(ctx context.Context, ids []string) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeByIDsWithEpisodes")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("anime.ids.count", len(ids))
	defer span.Finish()

	if len(ids) == 0 {
		return []*anime.Anime{}, nil
	}

	return a.Repository.FindByIDsWithEpisodes(spanCtx, ids)
}

// AnimeByIDs fetches multiple anime by their IDs in a single query
func (a *AnimeService) AnimeByIDs(ctx context.Context, ids []string) ([]*anime.Anime, error) {
	span, spanCtx := tracer.StartSpanFromContext(ctx, "AnimeByIDs")
	span.SetTag("service", "anime")
	span.SetTag("type", "service")
	span.SetTag("anime.ids.count", len(ids))
	defer span.Finish()

	if len(ids) == 0 {
		return []*anime.Anime{}, nil
	}

	return a.Repository.FindByIDs(spanCtx, ids)
}