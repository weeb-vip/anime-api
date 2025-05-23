package resolvers

import (
	"context"
	"encoding/json"
	metrics_lib "github.com/TempMee/go-metrics-lib"
	"github.com/weeb-vip/anime-api/graph/model"
	anime2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/metrics"
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
	if animeEntity.StartDate != nil {
		startDateTime, err := time.Parse("2006-01-02 15:04:05", *animeEntity.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &startDateTime
	}

	var endDate *time.Time
	if animeEntity.EndDate != nil {
		endDateTime, err := time.Parse("2006-01-02 15:04:05", *animeEntity.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &endDateTime
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
		EpisodeCount:  animeEntity.Episodes,
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

func transformAnimeToGraphQLWithEpisode(animeEntity anime2.AnimeWithNextEpisode) (*model.Anime, error) {

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
	if animeEntity.StartDate != nil {
		startDateTime, err := time.Parse("2006-01-02 15:04:05", *animeEntity.StartDate)
		if err != nil {
			return nil, err
		}
		startDate = &startDateTime
	}

	var endDate *time.Time
	if animeEntity.EndDate != nil {
		endDateTime, err := time.Parse("2006-01-02 15:04:05", *animeEntity.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &endDateTime
	}

	var nextEpisode *model.Episode

	if animeEntity.NextEpisode != nil {
		nextEpisodeEntity := animeEntity.NextEpisode
		nextEpisode = &model.Episode{
			ID:            nextEpisodeEntity.ID,
			AnimeID:       nextEpisodeEntity.AnimeID,
			EpisodeNumber: nextEpisodeEntity.Episode,
			TitleEn:       nextEpisodeEntity.TitleEn,
			TitleJp:       nextEpisodeEntity.TitleJp,
			AirDate:       nextEpisodeEntity.Aired,
			Synopsis:      nextEpisodeEntity.Synopsis,
			CreatedAt:     nextEpisodeEntity.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:     nextEpisodeEntity.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
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
		EpisodeCount:  animeEntity.Episodes,
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
		NextEpisode:   nextEpisode,
	}, nil
}

func AnimeByID(ctx context.Context, animeService anime.AnimeServiceImpl, id string) (*model.Anime, error) {

	startTime := time.Now()

	foundAnime, err := animeService.AnimeByID(ctx, id)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "AnimeByID",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
		})
		return nil, err
	}

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "AnimeByID",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return transformAnimeToGraphQL(*foundAnime)
}

func TopRatedAnime(ctx context.Context, animeService anime.AnimeServiceImpl, limit *int) ([]*model.Anime, error) {
	startTime := time.Now()

	if limit == nil {
		l := 10
		limit = &l
	}
	foundAnime, err := animeService.TopRatedAnime(ctx, *limit)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "TopRatedAnime",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
		})
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

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "TopRatedAnime",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return animes, nil
}

func MostPopularAnime(ctx context.Context, animeService anime.AnimeServiceImpl, limit *int) ([]*model.Anime, error) {
	startTime := time.Now()

	if limit == nil {
		l := 10
		limit = &l
	}
	foundAnime, err := animeService.MostPopularAnime(ctx, *limit)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "MostPopularAnime",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
		})
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

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "MostPopularAnime",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return animes, nil
}

func NewestAnime(ctx context.Context, animeService anime.AnimeServiceImpl, limit *int) ([]*model.Anime, error) {
	startTime := time.Now()

	if limit == nil {
		l := 10
		limit = &l
	}
	foundAnime, err := animeService.NewestAnime(ctx, *limit)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "NewestAnime",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
		})
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

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "NewestAnime",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return animes, nil
}

func CurrentlyAiring(ctx context.Context, animeService anime.AnimeServiceImpl, input *model.CurrentlyAiringInput) ([]*model.Anime, error) {
	startTime := time.Now()

	var foundAnime []*anime2.AnimeWithNextEpisode
	if input == nil {
		var err error
		foundAnime, err = animeService.AiringAnime(ctx, nil, nil, nil)
		if err != nil {
			_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
				Resolver: "CurrentlyAiring",
				Service:  "anime-api",
				Protocol: "graphql",
				Result:   metrics_lib.Error,
			})
			return nil, err
		}
	} else {
		var err error
		startDate := &input.StartDate
		foundAnime, err = animeService.AiringAnime(ctx, startDate, input.EndDate, input.DaysInFuture)
		if err != nil {
			_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
				Resolver: "CurrentlyAiring",
				Service:  "anime-api",
				Protocol: "graphql",
				Result:   metrics_lib.Error,
			})
			return nil, err
		}
	}

	var animes []*model.Anime
	for _, animeEntity := range foundAnime {
		anime, err := transformAnimeToGraphQLWithEpisode(*animeEntity)
		if err != nil {
			return nil, err
		}
		animes = append(animes, anime)
	}

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "CurrentlyAiring",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return animes, nil
}

func DBSearchAnime(ctx context.Context, animeService anime.AnimeServiceImpl, query string, page int, limit int) ([]*model.Anime, error) {
	startTime := time.Now()

	foundAnime, err := animeService.SearchedAnime(ctx, query, page, limit)
	if err != nil {
		_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
			Resolver: "DBSearchAnime",
			Service:  "anime-api",
			Protocol: "graphql",
			Result:   metrics_lib.Error,
		})
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

	_ = metrics.NewMetricsInstance().ResolverMetric(float64(time.Since(startTime).Milliseconds()), metrics_lib.ResolverMetricLabels{
		Resolver: "DBSearchAnime",
		Service:  "anime-api",
		Protocol: "graphql",
		Result:   metrics_lib.Success,
	})

	return animes, nil
}
