package episodes

import (
	"context"
	animeEpisode "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"time"
)

type AnimeEpisodeServiceImpl interface {
	GetEpisodesByAnimeID(ctx context.Context, animeID string) ([]*animeEpisode.AnimeEpisode, error)
	GetNextEpisode(ctx context.Context, animeID string) (*animeEpisode.AnimeEpisode, error)
}

type AnimeEpisodeService struct {
	Repository animeEpisode.AnimeEpisodeRepositoryImpl
}

func NewAnimeEpisodeService(repository animeEpisode.AnimeEpisodeRepositoryImpl) AnimeEpisodeServiceImpl {
	return &AnimeEpisodeService{
		Repository: repository,
	}
}

func (a *AnimeEpisodeService) GetEpisodesByAnimeID(ctx context.Context, animeID string) ([]*animeEpisode.AnimeEpisode, error) {
	return a.Repository.FindByAnimeID(ctx, animeID)
}

func (a *AnimeEpisodeService) GetNextEpisode(ctx context.Context, animeID string) (*animeEpisode.AnimeEpisode, error) {
	episodes, err := a.Repository.FindByAnimeID(ctx, animeID)
	if err != nil {
		return nil, err
	}

	if len(episodes) == 0 {
		return nil, nil
	}

	for _, episode := range episodes {
		// parse aired string to time.Time (we only want date, not time) ex: 2005-05-22 04:00:00
		var airedTime time.Time
		if episode.Aired != nil {
			airedTime = *episode.Aired
			if err != nil {
				return nil, err
			}

		} else {
			continue
		}

		if episode.Aired != nil && airedTime.After(time.Now()) {
			return episode, nil
		}
	}

	return nil, nil
}
