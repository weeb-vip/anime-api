package services

import (
	"testing"
	"time"

	"github.com/weeb-vip/anime-api/graph/model"
)

func TestProcessCurrentlyAiring(t *testing.T) {
	// Create test time
	now := time.Date(2023, 12, 15, 10, 0, 0, 0, time.UTC)

	// Create test episodes
	future1 := time.Date(2023, 12, 15, 11, 0, 0, 0, time.UTC) // 1 hour from now
	future2 := time.Date(2023, 12, 15, 12, 0, 0, 0, time.UTC) // 2 hours from now
	recent := time.Date(2023, 12, 15, 9, 45, 0, 0, time.UTC)  // 15 minutes ago (within 30min window)
	old := time.Date(2023, 12, 15, 8, 0, 0, 0, time.UTC)      // 2 hours ago (outside window)

	// Create test data with UTC broadcasts for simplicity
	animes := []*model.Anime{
		{
			ID:        "1",
			TitleEn:   stringPtr("Anime 1"),
			Duration:  stringPtr("24 min"),
			Broadcast: stringPtr("Fridays at 11:00 (UTC)"),
			Episodes: []*model.Episode{
				{
					ID:            "ep1",
					AnimeID:       stringPtr("1"),
					EpisodeNumber: intPtr(1),
					AirDate:       &future1,
				},
			},
		},
		{
			ID:        "2",
			TitleEn:   stringPtr("Anime 2"),
			Duration:  stringPtr("24 min"),
			Broadcast: stringPtr("Fridays at 12:00 (UTC)"),
			Episodes: []*model.Episode{
				{
					ID:            "ep2",
					AnimeID:       stringPtr("2"),
					EpisodeNumber: intPtr(1),
					AirDate:       &future2,
				},
			},
		},
		{
			ID:        "3",
			TitleEn:   stringPtr("Anime 3 - Recent"),
			Duration:  stringPtr("24 min"),
			Broadcast: stringPtr("Fridays at 09:45 (UTC)"),
			Episodes: []*model.Episode{
				{
					ID:            "ep3",
					AnimeID:       stringPtr("3"),
					EpisodeNumber: intPtr(1),
					AirDate:       &recent,
				},
			},
		},
		{
			ID:        "4",
			TitleEn:   stringPtr("Anime 4 - Old"),
			Duration:  stringPtr("24 min"),
			Broadcast: stringPtr("Fridays at 08:00 (UTC)"),
			Episodes: []*model.Episode{
				{
					ID:            "ep4",
					AnimeID:       stringPtr("4"),
					EpisodeNumber: intPtr(1),
					AirDate:       &old,
				},
			},
		},
	}

	// Test with limit of 10
	result := ProcessCurrentlyAiring(animes, 10, now)

	// Should return recently aired first (Anime 3), then future episodes (Anime 1, Anime 2)
	// Should exclude old episodes (Anime 4)
	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	// First should be recently aired
	if result[0].ID != "3" {
		t.Errorf("Expected first result to be Anime 3 (recently aired), got %s", result[0].ID)
	}

	// Second should be Anime 1 (next to air)
	if result[1].ID != "1" {
		t.Errorf("Expected second result to be Anime 1 (next future), got %s", result[1].ID)
	}

	// Third should be Anime 2
	if result[2].ID != "2" {
		t.Errorf("Expected third result to be Anime 2, got %s", result[2].ID)
	}

	// Test with limit of 2
	limitedResult := ProcessCurrentlyAiring(animes, 2, now)
	if len(limitedResult) != 2 {
		t.Errorf("Expected 2 results with limit, got %d", len(limitedResult))
	}
}

func TestParseAirTime(t *testing.T) {
	// Test JST timezone parsing
	airDate := time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC)
	broadcast := "Fridays at 01:30 (JST)"

	result := parseAirTime(&airDate, &broadcast)
	if result == nil {
		t.Error("Expected non-nil result")
		return
	}

	t.Logf("Input airDate: %v", airDate)
	t.Logf("Input broadcast: %s", broadcast)
	t.Logf("Result: %v", *result)

	// JST 01:30 should be UTC 16:30 (previous day)
	expected := time.Date(2023, 12, 14, 16, 30, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestCalculateCountdown(t *testing.T) {
	now := time.Date(2023, 12, 15, 10, 0, 0, 0, time.UTC)

	// Test future episode (1 hour away)
	future := time.Date(2023, 12, 15, 11, 0, 0, 0, time.UTC)
	broadcast := "Fridays at 11:00 (UTC)"

	countdown := calculateCountdown(&future, &broadcast, 24, now)
	if countdown != "1h" {
		t.Errorf("Expected '1h', got '%s'", countdown)
	}

	// Test just aired (30 minutes ago, past the 24 min episode duration)
	recent := time.Date(2023, 12, 15, 9, 30, 0, 0, time.UTC)
	recentBroadcast := "Fridays at 09:30 (UTC)"
	countdown = calculateCountdown(&recent, &recentBroadcast, 24, now)
	if countdown != "JUST AIRED" {
		t.Errorf("Expected 'JUST AIRED', got '%s'", countdown)
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}