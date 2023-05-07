package anime

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime"
)

type AnimeServiceImpl interface {
	AnimeByID(ctx context.Context, id string) (*anime.Anime, error)
	TopAnime(ctx context.Context, limit int) ([]*anime.Anime, error)
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
	return a.Repository.FindById(id)
}

func (a *AnimeService) TopAnime(ctx context.Context, limit int) ([]*anime.Anime, error) {
	return a.Repository.TopAnime(limit)
}
