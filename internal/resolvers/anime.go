package resolvers

import (
	"context"
	"encoding/json"
	"github.com/weeb-vip/anime-api/graph/model"
	anime2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime"
)

func transformAnimeToGraphQL(animeEntity anime2.Anime) (*model.Anime, error) {
	var studios []string
	if animeEntity.Studios != nil {
		err := json.Unmarshal([]byte(*animeEntity.Studios), &studios)
		if err != nil {
			return nil, err
		}
	}
	var tags []string
	if animeEntity.Genres != nil {
		err := json.Unmarshal([]byte(*animeEntity.Genres), &tags)
		if err != nil {
			return nil, err
		}
	}

	var titleSynonyms []string
	if animeEntity.TitleSynonyms != nil {
		err := json.Unmarshal([]byte(*animeEntity.TitleSynonyms), &titleSynonyms)
		if err != nil {
			return nil, err
		}
	}

	return &model.Anime{
		ID:            animeEntity.ID,
		TitleEn:       animeEntity.TitleEn,
		TitleJp:       animeEntity.TitleJp,
		TitleKanji:    animeEntity.TitleKanji,
		TitleRomaji:   animeEntity.TitleRomaji,
		TitleSynonyms: titleSynonyms,
		Description:   animeEntity.Synopsis,
		Episodes:      animeEntity.Episodes,
		Duration:      animeEntity.Duration,
		Studios:       studios,
		Tags:          tags,
		Rating:        animeEntity.Rating,
		AnimeStatus:   animeEntity.Status,
		ImageURL:      animeEntity.ImageURL,
		CreatedAt:     animeEntity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     animeEntity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func AnimeByID(ctx context.Context, animeService anime.AnimeServiceImpl, id string) (*model.Anime, error) {
	foundAnime, err := animeService.AnimeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transformAnimeToGraphQL(*foundAnime)
}

func TopAnime(ctx context.Context, animeService anime.AnimeServiceImpl, limit *int) ([]*model.Anime, error) {
	if limit == nil {
		l := 10
		limit = &l
	}
	foundAnime, err := animeService.TopAnime(ctx, *limit)
	if err != nil {
		return nil, err
	}

	var animes []*model.Anime
	for _, animeEntity := range foundAnime {
		anime, err := transformAnimeToGraphQL(*animeEntity)
		if err != nil {
			return nil, err
		}
		animes = append(animes, anime)
	}

	return animes, nil
}