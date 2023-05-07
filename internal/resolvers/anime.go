package resolvers

import (
	"context"
	"encoding/json"
	"github.com/weeb-vip/anime-api/graph/model"
	"github.com/weeb-vip/anime-api/internal/services/anime"
)

func AnimeByID(ctx context.Context, animeService anime.AnimeServiceImpl, id string) (*model.Anime, error) {
	foundAnime, err := animeService.AnimeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var studios []string
	if foundAnime.Studios != nil {
		err = json.Unmarshal([]byte(*foundAnime.Studios), &studios)
		if err != nil {
			return nil, err
		}
	}
	var tags []string
	if foundAnime.Genres != nil {
		err = json.Unmarshal([]byte(*foundAnime.Genres), &tags)
		if err != nil {
			return nil, err
		}
	}

	var titleSynonyms []string
	if foundAnime.TitleSynonyms != nil {
		err = json.Unmarshal([]byte(*foundAnime.TitleSynonyms), &titleSynonyms)
		if err != nil {
			return nil, err
		}
	}

	return &model.Anime{
		ID:            foundAnime.ID,
		TitleEn:       foundAnime.TitleEn,
		TitleJp:       foundAnime.TitleJp,
		TitleKanji:    foundAnime.TitleKanji,
		TitleRomaji:   foundAnime.TitleRomaji,
		TitleSynonyms: titleSynonyms,
		Description:   foundAnime.Synopsis,
		Episodes:      foundAnime.Episodes,
		Duration:      foundAnime.Duration,
		Studios:       studios,
		Tags:          tags,
		Rating:        foundAnime.Rating,
		AnimeStatus:   foundAnime.Status,
		ImageURL:      foundAnime.ImageURL,
		CreatedAt:     foundAnime.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     foundAnime.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil

}
