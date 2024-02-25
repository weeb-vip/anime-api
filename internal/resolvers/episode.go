package resolvers

import (
	"context"
	"github.com/weeb-vip/anime-api/graph/model"
	anime2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"github.com/weeb-vip/anime-api/internal/services/episodes"
	"time"
)

func transformEpisodeToGraphql(episodeEntity anime2.AnimeEpisode) (*model.Episode, error) {
	var aired *time.Time
	if episodeEntity.Aired != nil {
		airedTime, err := time.Parse("2006-01-02 15:04:05", *episodeEntity.Aired)
		if err != nil {
			return nil, err
		}
		aired = &airedTime
	}

	return &model.Episode{
		ID:            episodeEntity.ID,
		AnimeID:       episodeEntity.AnimeID,
		EpisodeNumber: episodeEntity.Episode,
		TitleEn:       episodeEntity.TitleEn,
		TitleJp:       episodeEntity.TitleJp,
		AirDate:       aired,
		Synopsis:      episodeEntity.Synopsis,
		CreatedAt:     episodeEntity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     episodeEntity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func EpisodesByAnimeID(ctx context.Context, animeEpisodeService episodes.AnimeEpisodeServiceImpl, animeID string) ([]*model.Episode, error) {
	episodeEntities, err := animeEpisodeService.GetEpisodesByAnimeID(ctx, animeID)
	if err != nil {
		return nil, err
	}

	episodes := make([]*model.Episode, len(episodeEntities))
	for i, episodeEntity := range episodeEntities {
		episode, err := transformEpisodeToGraphql(*episodeEntity)
		if err != nil {
			return nil, err
		}
		episodes[i] = episode
	}

	return episodes, nil
}
