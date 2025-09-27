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

	// Should return episodes from 30 minutes ago onwards in chronological order
	// Should include Anime 3 (15 min ago), Anime 1 (1h future), Anime 2 (2h future)
	// Should exclude Anime 4 (2h ago - outside 30min window)
	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	// Should be in chronological order: Anime 3 (earliest), then Anime 1, then Anime 2
	if result[0].ID != "3" {
		t.Errorf("Expected first result to be Anime 3 (earliest in chronological order), got %s", result[0].ID)
	}

	if result[1].ID != "1" {
		t.Errorf("Expected second result to be Anime 1 (next in chronological order), got %s", result[1].ID)
	}

	if result[2].ID != "2" {
		t.Errorf("Expected third result to be Anime 2 (last in chronological order), got %s", result[2].ID)
	}

	// Test that airTime is populated in UTC
	for _, anime := range result {
		if anime.NextEpisode != nil && anime.NextEpisode.AirTime != nil {
			// airTime should be in UTC
			if anime.NextEpisode.AirTime.Location() != time.UTC {
				t.Errorf("Expected airTime to be in UTC, got %s", anime.NextEpisode.AirTime.Location())
			}
			t.Logf("Anime %s: airDate=%v, airTime=%v (UTC)", anime.ID, anime.NextEpisode.AirDate, anime.NextEpisode.AirTime)
		}
	}

	// Test with limit of 2 - should get first 2 in chronological order
	limitedResult := ProcessCurrentlyAiring(animes, 2, now)
	if len(limitedResult) != 2 {
		t.Errorf("Expected 2 results with limit, got %d", len(limitedResult))
	}

	// Should be Anime 3 (earliest) and Anime 1 (next earliest)
	if limitedResult[0].ID != "3" {
		t.Errorf("Expected first limited result to be Anime 3, got %s", limitedResult[0].ID)
	}
	if limitedResult[1].ID != "1" {
		t.Errorf("Expected second limited result to be Anime 1, got %s", limitedResult[1].ID)
	}
}

func TestProcessCurrentlyAiring30MinuteWindow(t *testing.T) {
	// Test the 30-minute window behavior specifically
	now := time.Date(2023, 12, 15, 10, 0, 0, 0, time.UTC)

	// Create episodes at various times relative to "now"
	tooOld := time.Date(2023, 12, 15, 9, 20, 0, 0, time.UTC)     // 40 minutes ago (excluded)
	justInWindow := time.Date(2023, 12, 15, 9, 35, 0, 0, time.UTC) // 25 minutes ago (included)
	recent := time.Date(2023, 12, 15, 9, 55, 0, 0, time.UTC)    // 5 minutes ago (included)
	future1 := time.Date(2023, 12, 15, 10, 30, 0, 0, time.UTC)  // 30 minutes future (included)
	future2 := time.Date(2023, 12, 15, 11, 0, 0, 0, time.UTC)   // 1 hour future (included)

	animes := []*model.Anime{
		{
			ID: "too-old", TitleEn: stringPtr("Too Old"), Duration: stringPtr("24 min"), Broadcast: stringPtr("Fridays at 09:20 (UTC)"),
			Episodes: []*model.Episode{{ID: "ep1", AnimeID: stringPtr("too-old"), AirDate: &tooOld}},
		},
		{
			ID: "just-in-window", TitleEn: stringPtr("Just In Window"), Duration: stringPtr("24 min"), Broadcast: stringPtr("Fridays at 09:35 (UTC)"),
			Episodes: []*model.Episode{{ID: "ep2", AnimeID: stringPtr("just-in-window"), AirDate: &justInWindow}},
		},
		{
			ID: "recent", TitleEn: stringPtr("Recent"), Duration: stringPtr("24 min"), Broadcast: stringPtr("Fridays at 09:55 (UTC)"),
			Episodes: []*model.Episode{{ID: "ep3", AnimeID: stringPtr("recent"), AirDate: &recent}},
		},
		{
			ID: "future1", TitleEn: stringPtr("Future 1"), Duration: stringPtr("24 min"), Broadcast: stringPtr("Fridays at 10:30 (UTC)"),
			Episodes: []*model.Episode{{ID: "ep4", AnimeID: stringPtr("future1"), AirDate: &future1}},
		},
		{
			ID: "future2", TitleEn: stringPtr("Future 2"), Duration: stringPtr("24 min"), Broadcast: stringPtr("Fridays at 11:00 (UTC)"),
			Episodes: []*model.Episode{{ID: "ep5", AnimeID: stringPtr("future2"), AirDate: &future2}},
		},
	}

	result := ProcessCurrentlyAiring(animes, 10, now)

	// Should exclude "too-old" (40 min ago) but include the other 4
	expectedIDs := []string{"just-in-window", "recent", "future1", "future2"}
	if len(result) != len(expectedIDs) {
		t.Errorf("Expected %d results, got %d", len(expectedIDs), len(result))
	}

	// Should be in chronological order
	for i, expectedID := range expectedIDs {
		if i < len(result) && result[i].ID != expectedID {
			t.Errorf("Expected result[%d] to be %s, got %s", i, expectedID, result[i].ID)
		}
	}

	t.Logf("30-minute window test passed: excluded episodes older than 30 minutes, included others in chronological order")
}

func TestParseAirTime(t *testing.T) {
	// Test JST timezone parsing
	airDate := time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC)
	broadcast := "Fridays at 01:30 (JST)"

	result := ParseAirTime(&airDate, &broadcast)
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