package resolvers

import (
	"context"
	"encoding/json"
	"github.com/weeb-vip/anime-api/graph/model"
	anime2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/metrics"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"time"
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

	var licensors []string
	if animeEntity.Licensors != nil {
		err := json.Unmarshal([]byte(*animeEntity.Licensors), &licensors)
		if err != nil {
			return nil, err
		}
	}

	var startDate *time.Time
	if animeEntity.StartDate.Valid {
		startDate = &animeEntity.StartDate.Time
	}

	var endDate *time.Time
	if animeEntity.EndDate.Valid {
		endDate = &animeEntity.EndDate.Time
	}

	return &model.Anime{
		ID:            animeEntity.ID,
		Anidbid:       animeEntity.AnidbID,
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
		StartDate:     startDate,
		EndDate:       endDate,
		Broadcast:     animeEntity.Broadcast,
		Source:        animeEntity.Source,
		Licensors:     licensors,
		Ranking:       animeEntity.Ranking,
		CreatedAt:     animeEntity.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     animeEntity.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func AnimeByID(ctx context.Context, animeService anime.AnimeServiceImpl, id string) (*model.Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "AnimeByID")
	span.SetTag("id", id)
	span.SetTag("type", "resolver")
	defer span.Finish()
	startTime := time.Now()

	foundAnime, err := animeService.AnimeByID(ctx, id)
	if err != nil {
		metrics.ResolverHistrogramMetricError("AnimeByID", float64(time.Since(startTime).Milliseconds()))
		return nil, err
	}

	metrics.ResolverHistrogramMetricSuccess("AnimeByID", float64(time.Since(startTime).Milliseconds()))

	return transformAnimeToGraphQL(*foundAnime)
}

func TopRatedAnime(ctx context.Context, animeService anime.AnimeServiceImpl, limit *int) ([]*model.Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "TopRatedAnime")
	span.SetTag("type", "resolver")
	defer span.Finish()
	startTime := time.Now()

	if limit == nil {
		l := 10
		limit = &l
	}
	foundAnime, err := animeService.TopRatedAnime(ctx, *limit)
	if err != nil {
		metrics.ResolverHistrogramMetricError("TopRatedAnime", float64(time.Since(startTime).Milliseconds()))
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

	metrics.ResolverHistrogramMetricSuccess("TopRatedAnime", float64(time.Since(startTime).Milliseconds()))

	return animes, nil
}

func MostPopularAnime(ctx context.Context, animeService anime.AnimeServiceImpl, limit *int) ([]*model.Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "MostPopularAnime")
	span.SetTag("type", "resolver")
	defer span.Finish()
	startTime := time.Now()

	if limit == nil {
		l := 10
		limit = &l
	}
	foundAnime, err := animeService.MostPopularAnime(ctx, *limit)
	if err != nil {
		metrics.ResolverHistrogramMetricError("MostPopularAnime", float64(time.Since(startTime).Milliseconds()))
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

	metrics.ResolverHistrogramMetricSuccess("MostPopularAnime", float64(time.Since(startTime).Milliseconds()))

	return animes, nil
}

func NewestAnime(ctx context.Context, animeService anime.AnimeServiceImpl, limit *int) ([]*model.Anime, error) {
	span, _ := tracer.StartSpanFromContext(ctx, "NewestAnime")
	span.SetTag("type", "resolver")
	defer span.Finish()
	startTime := time.Now()

	if limit == nil {
		l := 10
		limit = &l
	}
	foundAnime, err := animeService.NewestAnime(ctx, *limit)
	if err != nil {
		metrics.ResolverHistrogramMetricError("NewestAnime", float64(time.Since(startTime).Milliseconds()))
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

	metrics.ResolverHistrogramMetricSuccess("NewestAnime", float64(time.Since(startTime).Milliseconds()))

	return animes, nil
}
