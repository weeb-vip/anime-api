package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	animerepo "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime"
)

// animeByIDStub answers AnimeByID and nothing else; the embedded interface
// keeps the compiler happy and panics loudly if the resolver strays.
type animeByIDStub struct {
	anime.AnimeServiceImpl
	byID map[string]*animerepo.Anime
	err  error
}

func (s *animeByIDStub) AnimeByID(_ context.Context, id string) (*animerepo.Anime, error) {
	if s.err != nil {
		return nil, s.err
	}
	if a, ok := s.byID[id]; ok {
		return a, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *animeByIDStub) AnimeByIDWithEpisodes(ctx context.Context, id string) (*animerepo.Anime, error) {
	return s.AnimeByID(ctx, id)
}

func strPtr(s string) *string { return &s }

func TestEntityResolver_FindAnimeByID(t *testing.T) {
	t.Run("resolves a reference another subgraph handed to the router", func(t *testing.T) {
		svc := &animeByIDStub{byID: map[string]*animerepo.Anime{
			"b977c2ab": {ID: "b977c2ab", TitleEn: strPtr("Koupen-chan"), UrlSlug: strPtr("koupen-chan")},
		}}
		r := &Resolver{AnimeService: svc}

		got, err := r.Entity().FindAnimeByID(context.Background(), "b977c2ab")

		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "b977c2ab", got.ID)
		assert.Equal(t, "Koupen-chan", *got.TitleEn)
	})

	t.Run("an unknown id is null, not an error", func(t *testing.T) {
		r := &Resolver{AnimeService: &animeByIDStub{byID: map[string]*animerepo.Anime{}}}

		got, err := r.Entity().FindAnimeByID(context.Background(), "missing")

		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("a real failure is still an error", func(t *testing.T) {
		r := &Resolver{AnimeService: &animeByIDStub{err: errors.New("db down")}}

		got, err := r.Entity().FindAnimeByID(context.Background(), "b977c2ab")

		require.Error(t, err)
		assert.Nil(t, got)
	})
}
