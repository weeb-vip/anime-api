// Package loaders batches the per-object fetches that a list of entities would
// otherwise issue one at a time.
//
// The case this exists for: UserAnime.anime. A watchlist comes from
// list-service, and the router asks this subgraph to resolve `anime` on every
// entry it returned. gqlgen resolves those concurrently but independently, so
// each one ran its own AnimeByID -- a page of 24 entries meant 24 database
// round trips, and the profile dashboard's thousand-entry lists meant a
// thousand. Anime.episodes had already been given a targeted N+1 guard for the
// seasonal query; this is the same disease on the watchlist path.
//
// AnimeByIDs and AnimeByIDsWithEpisodes were already written for exactly this
// and had no callers. All that was missing was something to gather the ids.
package loaders

import (
	"context"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"

	animerepo "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	animeservice "github.com/weeb-vip/anime-api/internal/services/anime"
)

type ctxKey string

const loadersKey ctxKey = "dataloaders"

// Loaders is one set per request. A loader caches within its own lifetime, so
// it must not outlive the request that created it -- two viewers would
// otherwise share a cached anime, and a long-lived loader would never see a
// row change.
type Loaders struct {
	// AnimeByID resolves without episodes.
	AnimeByID *dataloadgen.Loader[string, *animerepo.Anime]
	// AnimeWithEpisodesByID resolves with episodes preloaded, for callers whose
	// selection set asks for them. Kept separate rather than always preloading:
	// the join is far wider, and most callers do not want it.
	AnimeWithEpisodesByID *dataloadgen.Loader[string, *animerepo.Anime]
}

// byID reshapes a repository result into the order dataloadgen expects: one
// slot per requested key, in the order asked, nil where the row is missing.
func byID(ids []string, found []*animerepo.Anime, err error) ([]*animerepo.Anime, []error) {
	if err != nil {
		errs := make([]error, len(ids))
		for i := range errs {
			errs[i] = err
		}

		return make([]*animerepo.Anime, len(ids)), errs
	}

	index := make(map[string]*animerepo.Anime, len(found))
	for _, a := range found {
		if a != nil {
			index[a.ID] = a
		}
	}

	out := make([]*animerepo.Anime, len(ids))
	for i, id := range ids {
		// A missing id is a nil result, not an error: the field is nullable and
		// an entry pointing at a deleted anime should render as absent rather
		// than fail the whole query.
		out[i] = index[id]
	}

	return out, nil
}

// NewLoaders builds a request-scoped set.
func NewLoaders(animeService animeservice.AnimeServiceImpl) *Loaders {
	// 1ms is long enough to collect every sibling in a list -- gqlgen dispatches
	// them together -- and short enough to be invisible next to the query it
	// replaces.
	const wait = time.Millisecond

	return &Loaders{
		AnimeByID: dataloadgen.NewLoader(
			func(ctx context.Context, ids []string) ([]*animerepo.Anime, []error) {
				found, err := animeService.AnimeByIDs(ctx, ids)

				return byID(ids, found, err)
			},
			dataloadgen.WithWait(wait),
		),
		AnimeWithEpisodesByID: dataloadgen.NewLoader(
			func(ctx context.Context, ids []string) ([]*animerepo.Anime, []error) {
				found, err := animeService.AnimeByIDsWithEpisodes(ctx, ids)

				return byID(ids, found, err)
			},
			dataloadgen.WithWait(wait),
		),
	}
}

// Middleware attaches a fresh set of loaders to every request.
func Middleware(animeService animeservice.AnimeServiceImpl, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, NewLoaders(animeService))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// For returns the request's loaders, or nil when the caller reached a resolver
// without passing through Middleware. Callers fall back to the unbatched path
// rather than panicking, so a route that forgets the middleware is slow rather
// than broken.
func For(ctx context.Context) *Loaders {
	loaders, _ := ctx.Value(loadersKey).(*Loaders)

	return loaders
}
