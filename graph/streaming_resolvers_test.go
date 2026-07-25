package graph

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/weeb-vip/anime-api/graph/model"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_schedule"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_streaming_platform"
	"github.com/weeb-vip/anime-api/internal/db/repositories/episode_air_time"
)

// --- fake repositories implementing the AnimeSchedule integration interfaces ---

type fakeStreamingPlatformRepo struct {
	platforms []anime_streaming_platform.AnimeStreamingPlatform
	err       error
}

func (f *fakeStreamingPlatformRepo) FindByAnimeID(animeID string) ([]anime_streaming_platform.AnimeStreamingPlatform, error) {
	return f.platforms, f.err
}

type fakeScheduleRepo struct {
	schedule *anime_schedule.AnimeSchedule
	err      error
}

func (f *fakeScheduleRepo) FindByAnimeID(animeID string) (*anime_schedule.AnimeSchedule, error) {
	return f.schedule, f.err
}

type fakeEpisodeAirTimeRepo struct {
	airTimes []episode_air_time.EpisodeAirTime
	err      error
}

func (f *fakeEpisodeAirTimeRepo) FindByAnimeID(animeID string) ([]episode_air_time.EpisodeAirTime, error) {
	return f.airTimes, f.err
}

func (f *fakeEpisodeAirTimeRepo) FindByAnimeIDAndEpisode(animeID string, episodeNumber int) ([]episode_air_time.EpisodeAirTime, error) {
	return f.airTimes, f.err
}

func (f *fakeEpisodeAirTimeRepo) FindSubTimeByAnimeIDAndEpisode(animeID string, episodeNumber int) (*episode_air_time.EpisodeAirTime, error) {
	if f.err != nil {
		return nil, f.err
	}
	if len(f.airTimes) == 0 {
		return nil, nil
	}
	return &f.airTimes[0], nil
}

func TestAnimeResolver_StreamingPlatforms(t *testing.T) {
	ctx := context.Background()
	obj := &model.Anime{ID: "anime-1"}

	t.Run("nil repository returns nil without error", func(t *testing.T) {
		r := &Resolver{}
		got, err := r.Anime().StreamingPlatforms(ctx, obj)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("repository error degrades to nil", func(t *testing.T) {
		r := &Resolver{AnimeStreamingPlatformRepository: &fakeStreamingPlatformRepo{err: errors.New("db down")}}
		got, err := r.Anime().StreamingPlatforms(ctx, obj)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("maps platforms preserving order and optional name", func(t *testing.T) {
		crName := "Crunchyroll"
		r := &Resolver{AnimeStreamingPlatformRepository: &fakeStreamingPlatformRepo{
			platforms: []anime_streaming_platform.AnimeStreamingPlatform{
				{Platform: "crunchyroll", Name: &crName, URL: "https://crunchyroll.com/x"},
				{Platform: "netflix", Name: nil, URL: "https://netflix.com/y"},
			},
		}}
		got, err := r.Anime().StreamingPlatforms(ctx, obj)
		require.NoError(t, err)
		require.Len(t, got, 2)

		assert.Equal(t, "crunchyroll", got[0].Platform)
		require.NotNil(t, got[0].Name)
		assert.Equal(t, "Crunchyroll", *got[0].Name)
		assert.Equal(t, "https://crunchyroll.com/x", got[0].URL)

		assert.Equal(t, "netflix", got[1].Platform)
		assert.Nil(t, got[1].Name)
		assert.Equal(t, "https://netflix.com/y", got[1].URL)
	})

	t.Run("empty result returns nil", func(t *testing.T) {
		r := &Resolver{AnimeStreamingPlatformRepository: &fakeStreamingPlatformRepo{platforms: nil}}
		got, err := r.Anime().StreamingPlatforms(ctx, obj)
		require.NoError(t, err)
		assert.Nil(t, got)
	})
}

func TestAnimeResolver_ScheduleInfo(t *testing.T) {
	ctx := context.Background()
	obj := &model.Anime{ID: "anime-1"}

	t.Run("nil repository returns nil without error", func(t *testing.T) {
		r := &Resolver{}
		got, err := r.Anime().ScheduleInfo(ctx, obj)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("not found degrades to nil", func(t *testing.T) {
		r := &Resolver{AnimeScheduleRepository: &fakeScheduleRepo{err: errors.New("record not found")}}
		got, err := r.Anime().ScheduleInfo(ctx, obj)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("maps schedule fields", func(t *testing.T) {
		jpn := time.Date(2026, 1, 2, 15, 0, 0, 0, time.UTC)
		notes := "delayed one week"
		delayed := "2026-01-09"
		r := &Resolver{AnimeScheduleRepository: &fakeScheduleRepo{schedule: &anime_schedule.AnimeSchedule{
			AnimeID:          "anime-1",
			JpnTime:          &jpn,
			Notes:            &notes,
			DelayedTimetable: &delayed,
		}}}
		got, err := r.Anime().ScheduleInfo(ctx, obj)
		require.NoError(t, err)
		require.NotNil(t, got)
		require.NotNil(t, got.JpnTime)
		assert.Equal(t, jpn, *got.JpnTime)
		require.NotNil(t, got.Notes)
		assert.Equal(t, "delayed one week", *got.Notes)
		require.NotNil(t, got.DelayedTimetable)
		assert.Equal(t, "2026-01-09", *got.DelayedTimetable)
		// Untouched optional fields stay nil.
		assert.Nil(t, got.SubTime)
		assert.Nil(t, got.DubTime)
	})
}

func TestEpisodeResolver_AirTimes(t *testing.T) {
	ctx := context.Background()
	animeID := "anime-1"
	epNum := 3
	obj := &model.Episode{AnimeID: &animeID, EpisodeNumber: &epNum}

	t.Run("nil repository returns nil", func(t *testing.T) {
		r := &Resolver{}
		got, err := r.Episode().AirTimes(ctx, obj)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("nil animeID returns nil", func(t *testing.T) {
		r := &Resolver{EpisodeAirTimeRepository: &fakeEpisodeAirTimeRepo{}}
		got, err := r.Episode().AirTimes(ctx, &model.Episode{EpisodeNumber: &epNum})
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("nil episodeNumber returns nil", func(t *testing.T) {
		r := &Resolver{EpisodeAirTimeRepository: &fakeEpisodeAirTimeRepo{}}
		got, err := r.Episode().AirTimes(ctx, &model.Episode{AnimeID: &animeID})
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("repository error degrades to nil", func(t *testing.T) {
		r := &Resolver{EpisodeAirTimeRepository: &fakeEpisodeAirTimeRepo{err: errors.New("db down")}}
		got, err := r.Episode().AirTimes(ctx, obj)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("parses streams JSON and maps air type", func(t *testing.T) {
		when := time.Date(2026, 1, 2, 15, 0, 0, 0, time.UTC)
		streams := `[{"platform":"crunchyroll","name":"Crunchyroll","url":"https://cr.com/x"},{"platform":"netflix","name":"","url":"https://nf.com/y"}]`
		r := &Resolver{EpisodeAirTimeRepository: &fakeEpisodeAirTimeRepo{airTimes: []episode_air_time.EpisodeAirTime{
			{AirType: "sub", AirDatetime: when, StreamsJSON: &streams},
		}}}
		got, err := r.Episode().AirTimes(ctx, obj)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, model.AirType("sub"), got[0].AirType)
		assert.Equal(t, when, got[0].AirDatetime)
		require.Len(t, got[0].Streams, 2)
		assert.Equal(t, "crunchyroll", got[0].Streams[0].Platform)
		require.NotNil(t, got[0].Streams[0].Name)
		assert.Equal(t, "Crunchyroll", *got[0].Streams[0].Name)
		assert.Equal(t, "https://cr.com/x", got[0].Streams[0].URL)
		assert.Equal(t, "netflix", got[0].Streams[1].Platform)
	})

	t.Run("nil streams JSON yields entry with no streams", func(t *testing.T) {
		when := time.Date(2026, 1, 2, 15, 0, 0, 0, time.UTC)
		r := &Resolver{EpisodeAirTimeRepository: &fakeEpisodeAirTimeRepo{airTimes: []episode_air_time.EpisodeAirTime{
			{AirType: "raw", AirDatetime: when, StreamsJSON: nil},
		}}}
		got, err := r.Episode().AirTimes(ctx, obj)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, model.AirType("raw"), got[0].AirType)
		assert.Nil(t, got[0].Streams)
	})

	t.Run("invalid streams JSON is ignored gracefully", func(t *testing.T) {
		when := time.Date(2026, 1, 2, 15, 0, 0, 0, time.UTC)
		bad := `{not valid json`
		r := &Resolver{EpisodeAirTimeRepository: &fakeEpisodeAirTimeRepo{airTimes: []episode_air_time.EpisodeAirTime{
			{AirType: "dub", AirDatetime: when, StreamsJSON: &bad},
		}}}
		got, err := r.Episode().AirTimes(ctx, obj)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, model.AirType("dub"), got[0].AirType)
		assert.Nil(t, got[0].Streams)
	})
}
