package resolvers

import (
	"context"
	"reflect"
	"testing"
	"time"

	anime_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	anime_episode_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	anime_season_repo "github.com/weeb-vip/anime-api/internal/db/repositories/anime_season"
	"go.uber.org/mock/gomock"
)

// MockAnimeService implements the AnimeServiceImpl interface for testing
type MockAnimeService struct {
	ctrl     *gomock.Controller
	recorder *MockAnimeServiceMockRecorder
}

type MockAnimeServiceMockRecorder struct {
	mock *MockAnimeService
}

func NewMockAnimeService(ctrl *gomock.Controller) *MockAnimeService {
	mock := &MockAnimeService{ctrl: ctrl}
	mock.recorder = &MockAnimeServiceMockRecorder{mock}
	return mock
}

func (m *MockAnimeService) EXPECT() *MockAnimeServiceMockRecorder {
	return m.recorder
}

func (m *MockAnimeService) AnimeByID(ctx context.Context, id string) (*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeByID", ctx, id)
	ret0, _ := ret[0].(*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeByID(ctx, id interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeByID", reflect.TypeOf((*MockAnimeService)(nil).AnimeByID), ctx, id)
}

func (m *MockAnimeService) AnimeByIDWithEpisodes(ctx context.Context, id string) (*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeByIDWithEpisodes", ctx, id)
	ret0, _ := ret[0].(*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeByIDWithEpisodes(ctx, id interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeByIDWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).AnimeByIDWithEpisodes), ctx, id)
}

func (m *MockAnimeService) AnimeByIDs(ctx context.Context, ids []string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeByIDs", ctx, ids)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeByIDs(ctx, ids interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeByIDs", reflect.TypeOf((*MockAnimeService)(nil).AnimeByIDs), ctx, ids)
}

func (m *MockAnimeService) AnimeByIDsWithEpisodes(ctx context.Context, ids []string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeByIDsWithEpisodes", ctx, ids)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeByIDsWithEpisodes(ctx, ids interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeByIDsWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).AnimeByIDsWithEpisodes), ctx, ids)
}

func (m *MockAnimeService) TopRatedAnime(ctx context.Context, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TopRatedAnime", ctx, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) TopRatedAnime(ctx, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "TopRatedAnime", reflect.TypeOf((*MockAnimeService)(nil).TopRatedAnime), ctx, limit)
}

func (m *MockAnimeService) TopRatedAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TopRatedAnimeWithEpisodes", ctx, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) TopRatedAnimeWithEpisodes(ctx, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "TopRatedAnimeWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).TopRatedAnimeWithEpisodes), ctx, limit)
}

func (m *MockAnimeService) MostPopularAnime(ctx context.Context, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MostPopularAnime", ctx, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) MostPopularAnime(ctx, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "MostPopularAnime", reflect.TypeOf((*MockAnimeService)(nil).MostPopularAnime), ctx, limit)
}

func (m *MockAnimeService) MostPopularAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MostPopularAnimeWithEpisodes", ctx, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) MostPopularAnimeWithEpisodes(ctx, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "MostPopularAnimeWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).MostPopularAnimeWithEpisodes), ctx, limit)
}

func (m *MockAnimeService) NewestAnime(ctx context.Context, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewestAnime", ctx, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) NewestAnime(ctx, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "NewestAnime", reflect.TypeOf((*MockAnimeService)(nil).NewestAnime), ctx, limit)
}

func (m *MockAnimeService) NewestAnimeWithEpisodes(ctx context.Context, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewestAnimeWithEpisodes", ctx, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) NewestAnimeWithEpisodes(ctx, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "NewestAnimeWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).NewestAnimeWithEpisodes), ctx, limit)
}

func (m *MockAnimeService) AiringAnime(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime_repo.AnimeWithNextEpisode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AiringAnime", ctx, startDate, endDate, days)
	ret0, _ := ret[0].([]*anime_repo.AnimeWithNextEpisode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AiringAnime(ctx, startDate, endDate, days interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AiringAnime", reflect.TypeOf((*MockAnimeService)(nil).AiringAnime), ctx, startDate, endDate, days)
}

func (m *MockAnimeService) AiringAnimeWithEpisodes(ctx context.Context, startDate *time.Time, endDate *time.Time, days *int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AiringAnimeWithEpisodes", ctx, startDate, endDate, days)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AiringAnimeWithEpisodes(ctx, startDate, endDate, days interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AiringAnimeWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).AiringAnimeWithEpisodes), ctx, startDate, endDate, days)
}

func (m *MockAnimeService) SearchedAnime(ctx context.Context, query string, page int, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchedAnime", ctx, query, page, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) SearchedAnime(ctx, query, page, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "SearchedAnime", reflect.TypeOf((*MockAnimeService)(nil).SearchedAnime), ctx, query, page, limit)
}

func (m *MockAnimeService) SearchedAnimeWithEpisodes(ctx context.Context, query string, page int, limit int) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchedAnimeWithEpisodes", ctx, query, page, limit)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) SearchedAnimeWithEpisodes(ctx, query, page, limit interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "SearchedAnimeWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).SearchedAnimeWithEpisodes), ctx, query, page, limit)
}

func (m *MockAnimeService) AnimeBySeasonWithEpisodes(ctx context.Context, season string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeBySeasonWithEpisodes", ctx, season)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeBySeasonWithEpisodes(ctx, season interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeBySeasonWithEpisodes", reflect.TypeOf((*MockAnimeService)(nil).AnimeBySeasonWithEpisodes), ctx, season)
}

func (m *MockAnimeService) AnimeBySeasonWithEpisodesOptimized(ctx context.Context, season string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeBySeasonWithEpisodesOptimized", ctx, season)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeBySeasonWithEpisodesOptimized(ctx, season interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeBySeasonWithEpisodesOptimized", reflect.TypeOf((*MockAnimeService)(nil).AnimeBySeasonWithEpisodesOptimized), ctx, season)
}

func (m *MockAnimeService) AnimeBySeasonAnimeOnlyOptimized(ctx context.Context, season string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeBySeasonAnimeOnlyOptimized", ctx, season)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeBySeasonAnimeOnlyOptimized(ctx, season interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeBySeasonAnimeOnlyOptimized", reflect.TypeOf((*MockAnimeService)(nil).AnimeBySeasonAnimeOnlyOptimized), ctx, season)
}

func (m *MockAnimeService) AnimeBySeasonWithIndexHints(ctx context.Context, season string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeBySeasonWithIndexHints", ctx, season)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeBySeasonWithIndexHints(ctx, season interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeBySeasonWithIndexHints", reflect.TypeOf((*MockAnimeService)(nil).AnimeBySeasonWithIndexHints), ctx, season)
}

func (m *MockAnimeService) AnimeBySeasonBatched(ctx context.Context, season string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeBySeasonBatched", ctx, season)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeBySeasonBatched(ctx, season interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeBySeasonBatched", reflect.TypeOf((*MockAnimeService)(nil).AnimeBySeasonBatched), ctx, season)
}

func (m *MockAnimeService) AnimeBySeasonOptimized(ctx context.Context, season string) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeBySeasonOptimized", ctx, season)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeBySeasonOptimized(ctx, season interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeBySeasonOptimized", reflect.TypeOf((*MockAnimeService)(nil).AnimeBySeasonOptimized), ctx, season)
}

func (m *MockAnimeService) AnimeBySeasonWithFieldSelection(ctx context.Context, season string, fields *anime_repo.FieldSelection) ([]*anime_repo.Anime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnimeBySeasonWithFieldSelection", ctx, season, fields)
	ret0, _ := ret[0].([]*anime_repo.Anime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeServiceMockRecorder) AnimeBySeasonWithFieldSelection(ctx, season, fields interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "AnimeBySeasonWithFieldSelection", reflect.TypeOf((*MockAnimeService)(nil).AnimeBySeasonWithFieldSelection), ctx, season, fields)
}

// MockAnimeSeasonService implements the AnimeSeasonServiceImpl interface for testing
type MockAnimeSeasonService struct {
	ctrl     *gomock.Controller
	recorder *MockAnimeSeasonServiceMockRecorder
}

type MockAnimeSeasonServiceMockRecorder struct {
	mock *MockAnimeSeasonService
}

func NewMockAnimeSeasonService(ctrl *gomock.Controller) *MockAnimeSeasonService {
	mock := &MockAnimeSeasonService{ctrl: ctrl}
	mock.recorder = &MockAnimeSeasonServiceMockRecorder{mock}
	return mock
}

func (m *MockAnimeSeasonService) EXPECT() *MockAnimeSeasonServiceMockRecorder {
	return m.recorder
}

// Note: These methods are part of the interface but not used in our test
func (m *MockAnimeSeasonService) FindByAnimeID(ctx context.Context, animeID string) ([]*anime_season_repo.AnimeSeason, error) {
	// Implementation not needed for this test
	return nil, nil
}

func (m *MockAnimeSeasonService) FindBySeason(ctx context.Context, season string) ([]*anime_season_repo.AnimeSeason, error) {
	// Implementation not needed for this test
	return nil, nil
}

func (m *MockAnimeSeasonService) Create(ctx context.Context, animeSeason *anime_season_repo.AnimeSeason) error {
	// Implementation not needed for this test
	return nil
}

func (m *MockAnimeSeasonService) Update(ctx context.Context, animeSeason *anime_season_repo.AnimeSeason) error {
	// Implementation not needed for this test
	return nil
}

func (m *MockAnimeSeasonService) Delete(ctx context.Context, id string) error {
	// Implementation not needed for this test
	return nil
}

// MockAnimeEpisodeService implements the AnimeEpisodeServiceImpl interface for testing
type MockAnimeEpisodeService struct {
	ctrl     *gomock.Controller
	recorder *MockAnimeEpisodeServiceMockRecorder
}

type MockAnimeEpisodeServiceMockRecorder struct {
	mock *MockAnimeEpisodeService
}

func NewMockAnimeEpisodeService(ctrl *gomock.Controller) *MockAnimeEpisodeService {
	mock := &MockAnimeEpisodeService{ctrl: ctrl}
	mock.recorder = &MockAnimeEpisodeServiceMockRecorder{mock}
	return mock
}

func (m *MockAnimeEpisodeService) EXPECT() *MockAnimeEpisodeServiceMockRecorder {
	return m.recorder
}

func (m *MockAnimeEpisodeService) GetEpisodesByAnimeID(ctx context.Context, animeID string) ([]*anime_episode_repo.AnimeEpisode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEpisodesByAnimeID", ctx, animeID)
	ret0, _ := ret[0].([]*anime_episode_repo.AnimeEpisode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeEpisodeServiceMockRecorder) GetEpisodesByAnimeID(ctx, animeID interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "GetEpisodesByAnimeID", reflect.TypeOf((*MockAnimeEpisodeService)(nil).GetEpisodesByAnimeID), ctx, animeID)
}

func (m *MockAnimeEpisodeService) GetNextEpisode(ctx context.Context, animeID string) (*anime_episode_repo.AnimeEpisode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextEpisode", ctx, animeID)
	ret0, _ := ret[0].(*anime_episode_repo.AnimeEpisode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (c *MockAnimeEpisodeServiceMockRecorder) GetNextEpisode(ctx, animeID interface{}) *gomock.Call {
	c.mock.ctrl.T.Helper()
	return c.mock.ctrl.RecordCallWithMethodType(c.mock, "GetNextEpisode", reflect.TypeOf((*MockAnimeEpisodeService)(nil).GetNextEpisode), ctx, animeID)
}

func TestAnimeBySeasonsNoAdditionalEpisodeQueries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := NewMockAnimeService(ctrl)
	mockAnimeSeasonService := NewMockAnimeSeasonService(ctrl)
	mockAnimeEpisodeService := NewMockAnimeEpisodeService(ctrl)

	ctx := context.Background()
	season := "SPRING_2024"

	// Create test data with preloaded episodes
	now := time.Now()
	testAnime := &anime_repo.Anime{
		ID:       "anime-1",
		TitleEn:  stringPtr("Test Anime"),
		Episodes: intPtr(12),
		AnimeEpisodes: []*anime_episode_repo.AnimeEpisode{
			{
				ID:        "ep-1",
				AnimeID:   stringPtr("anime-1"),
				Episode:   intPtr(1),
				TitleEn:   stringPtr("Episode 1"),
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "ep-2",
				AnimeID:   stringPtr("anime-1"),
				Episode:   intPtr(2),
				TitleEn:   stringPtr("Episode 2"),
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Expect only ONE call to the optimized AnimeBySeasonOptimized method (without episodes)
	mockAnimeService.EXPECT().
		AnimeBySeasonOptimized(ctx, season).
		Return([]*anime_repo.Anime{testAnime}, nil).
		Times(1)

	// CRITICAL: We should NOT expect any calls to GetEpisodesByAnimeID
	// If this expectation fails, it means additional episode queries are being made
	mockAnimeEpisodeService.EXPECT().
		GetEpisodesByAnimeID(gomock.Any(), gomock.Any()).
		Times(0) // Zero calls expected

	// Execute the resolver
	result, err := AnimeBySeasons(ctx, mockAnimeSeasonService, mockAnimeService, season)

	// Verify results
	if err != nil {
		t.Fatalf("AnimeBySeasons returned error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 anime, got %d", len(result))
	}

	// Verify that episodes are properly included
	anime := result[0]
	if anime.ID != "anime-1" {
		t.Errorf("Expected anime ID 'anime-1', got '%s'", anime.ID)
	}

	if len(anime.Episodes) != 2 {
		t.Errorf("Expected 2 episodes, got %d", len(anime.Episodes))
	}

	// Verify episode data is correct
	if anime.Episodes[0].ID != "ep-1" {
		t.Errorf("Expected first episode ID 'ep-1', got '%s'", anime.Episodes[0].ID)
	}

	if anime.Episodes[1].ID != "ep-2" {
		t.Errorf("Expected second episode ID 'ep-2', got '%s'", anime.Episodes[1].ID)
	}
}

func TestAnimeBySeasonsWithEmptyEpisodesNoAdditionalQueries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAnimeService := NewMockAnimeService(ctrl)
	mockAnimeSeasonService := NewMockAnimeSeasonService(ctrl)
	mockAnimeEpisodeService := NewMockAnimeEpisodeService(ctrl)

	ctx := context.Background()
	season := "SUMMER_2024"

	// Create test data with NO episodes (empty slice)
	now := time.Now()
	testAnime := &anime_repo.Anime{
		ID:            "anime-2",
		TitleEn:       stringPtr("Test Anime No Episodes"),
		Episodes:      intPtr(0),
		AnimeEpisodes: []*anime_episode_repo.AnimeEpisode{}, // Empty slice
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Expect only ONE call to the optimized AnimeBySeasonOptimized method (without episodes)
	mockAnimeService.EXPECT().
		AnimeBySeasonOptimized(ctx, season).
		Return([]*anime_repo.Anime{testAnime}, nil).
		Times(1)

	// CRITICAL: We should NOT expect any calls to GetEpisodesByAnimeID
	// Even with empty episodes, we should not trigger additional queries
	mockAnimeEpisodeService.EXPECT().
		GetEpisodesByAnimeID(gomock.Any(), gomock.Any()).
		Times(0) // Zero calls expected

	// Execute the resolver
	result, err := AnimeBySeasons(ctx, mockAnimeSeasonService, mockAnimeService, season)

	// Verify results
	if err != nil {
		t.Fatalf("AnimeBySeasons returned error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 1 anime, got %d", len(result))
	}

	// Verify that episodes are empty but not nil
	anime := result[0]
	if anime.ID != "anime-2" {
		t.Errorf("Expected anime ID 'anime-2', got '%s'", anime.ID)
	}

	if anime.Episodes == nil {
		t.Error("Expected episodes to be empty slice, got nil")
	}

	if len(anime.Episodes) != 0 {
		t.Errorf("Expected 0 episodes, got %d", len(anime.Episodes))
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}