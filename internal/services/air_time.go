package services

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/weeb-vip/anime-api/graph/model"
)

// Episode represents an anime episode for air time calculations
type Episode struct {
	ID            string
	AnimeID       string
	EpisodeNumber *int
	TitleEn       *string
	TitleJp       *string
	AirDate       *time.Time
	Synopsis      *string
}

// NextEpisodeResult represents the result of finding the next episode
type NextEpisodeResult struct {
	Episode Episode
	AirTime time.Time
}

// AirTimeDisplayInfo represents the display information for air time
type AirTimeDisplayInfo struct {
	Show    bool
	Text    string
	Variant string // countdown, scheduled, aired, airing
}

// ParseAirTime parses broadcast time and creates accurate air time in UTC
// Based on the frontend airTimeUtils.js parseAirTime function
// Returns time in UTC timezone for consistent API responses
func ParseAirTime(airDate *time.Time, broadcast *string) *time.Time {
	if airDate == nil || broadcast == nil {
		return airDate
	}

	// Parse broadcast time (e.g., "Wednesdays at 01:29 (JST)")
	broadcastStr := *broadcast
	if !strings.Contains(broadcastStr, ":") {
		return airDate
	}

	// Extract time and timezone from broadcast string
	var hours, minutes int
	var timezone string

	// Look for time pattern HH:MM using regex-like approach
	// Find "at HH:MM" pattern
	atIndex := strings.Index(broadcastStr, " at ")
	if atIndex != -1 {
		timeStr := broadcastStr[atIndex+4:] // Skip " at "
		// Find the time pattern before any parentheses
		parenIndex := strings.Index(timeStr, " (")
		if parenIndex != -1 {
			timeStr = timeStr[:parenIndex]
		}

		n, err := fmt.Sscanf(timeStr, "%d:%d", &hours, &minutes)
		if err != nil || n != 2 {
			// If we can't parse the time, just return the original air date
			return airDate
		}
	} else {
		// Try to find just HH:MM pattern
		n, err := fmt.Sscanf(broadcastStr, "%d:%d", &hours, &minutes)
		if err != nil || n != 2 {
			// If we can't parse the time, just return the original air date
			return airDate
		}
	}

	// Extract timezone (JST, UTC, etc.)
	if strings.Contains(broadcastStr, "(JST)") {
		timezone = "JST"
	} else if strings.Contains(broadcastStr, "(UTC)") {
		timezone = "UTC"
	} else {
		timezone = "JST" // Default to JST for Japanese anime
	}

	// Create a date object from the air date
	parsedAirDate := *airDate

	if timezone == "JST" {
		// JST is UTC+9, so to get UTC time we subtract 9 hours
		utcHours := hours - 9
		utcDay := parsedAirDate.Day()
		utcMonth := parsedAirDate.Month()
		utcYear := parsedAirDate.Year()

		if utcHours < 0 {
			// Previous day case
			utcHours += 24
			// Go back one day
			tempDate := time.Date(utcYear, utcMonth, utcDay, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)
			utcYear = tempDate.Year()
			utcMonth = tempDate.Month()
			utcDay = tempDate.Day()
		}

		result := time.Date(utcYear, utcMonth, utcDay, utcHours, minutes, 0, 0, time.UTC)
		return &result
	} else {
		// For other timezones, assume UTC
		result := time.Date(
			parsedAirDate.Year(),
			parsedAirDate.Month(),
			parsedAirDate.Day(),
			hours,
			minutes,
			0, 0, time.UTC,
		)
		return &result
	}
}

// parseDurationToMinutes parses duration string to minutes
func parseDurationToMinutes(duration *string) int {
	if duration == nil {
		return 24 // Default to 24 minutes for typical anime episode
	}

	durationStr := *duration
	// Extract number from duration string (e.g., "24 min per episode", "24 min")
	var minutes int
	if _, err := fmt.Sscanf(durationStr, "%d", &minutes); err != nil {
		return 24 // Default fallback
	}

	return minutes
}

// isCurrentlyAiring checks if the anime is currently airing
func isCurrentlyAiring(airDate *time.Time, broadcast *string, durationMinutes int, currentTime time.Time) bool {
	if airDate == nil {
		return false
	}

	airTime := ParseAirTime(airDate, broadcast)
	if airTime == nil {
		return false
	}

	airStartMs := airTime.UnixMilli()
	currentMs := currentTime.UnixMilli()
	episodeDurationMs := int64(durationMinutes * 60 * 1000)
	airEndMs := airStartMs + episodeDurationMs

	// Currently airing if current time is between start and end time
	return currentMs >= airStartMs && currentMs <= airEndMs
}

// hasAlreadyAired checks if the anime has already aired
func hasAlreadyAired(airDate *time.Time, broadcast *string, durationMinutes int, currentTime time.Time) bool {
	if airDate == nil {
		return false
	}

	airTime := ParseAirTime(airDate, broadcast)
	if airTime == nil {
		return false
	}

	currentMs := currentTime.UnixMilli()
	airStartMs := airTime.UnixMilli()
	episodeDurationMs := int64(durationMinutes * 60 * 1000)
	airEndMs := airStartMs + episodeDurationMs
	sevenDaysMs := int64(7 * 24 * 60 * 60 * 1000)

	// Show "already aired" if episode has finished and it's within the last 7 days
	return currentMs > airEndMs && (currentMs-airEndMs) <= sevenDaysMs
}

// calculateCountdown calculates countdown string or current airing status
func calculateCountdown(airDate *time.Time, broadcast *string, durationMinutes int, currentTime time.Time) string {
	if airDate == nil {
		return ""
	}

	airTime := ParseAirTime(airDate, broadcast)
	if airTime == nil {
		return ""
	}

	// Check if currently airing
	if isCurrentlyAiring(airDate, broadcast, durationMinutes, currentTime) {
		airStartMs := airTime.UnixMilli()
		currentMs := currentTime.UnixMilli()
		episodeDurationMs := int64(durationMinutes * 60 * 1000)
		remainingMs := airStartMs + episodeDurationMs - currentMs
		remainingMinutes := remainingMs / (1000 * 60)

		if remainingMinutes < 60 {
			if remainingMinutes <= 0 {
				return "AIRING NOW"
			}
			return fmt.Sprintf("%dm left", remainingMinutes)
		}
		return "AIRING NOW"
	}

	diffMs := airTime.UnixMilli() - currentTime.UnixMilli()
	dayMs := int64(24 * 60 * 60 * 1000)

	// Upcoming today (within next 24h)
	if diffMs > 0 && diffMs <= dayMs {
		diffMinutes := diffMs / (1000 * 60)
		if diffMinutes < 60 {
			return fmt.Sprintf("%dm", diffMinutes)
		}
		diffHours := diffMinutes / 60
		return fmt.Sprintf("%dh", diffHours)
	}

	// Already started (but not currently airing) => just aired
	if diffMs <= 0 {
		return "JUST AIRED"
	}

	// Not today (>24h away)
	return ""
}

// findNextEpisode finds the next episode from an episodes array
func findNextEpisode(episodes []*model.Episode, broadcast *string, currentTime time.Time) *NextEpisodeResult {
	if len(episodes) == 0 {
		return nil
	}

	for _, episode := range episodes {
		if episode.AirDate != nil {
			airTime := ParseAirTime(episode.AirDate, broadcast)
			if airTime != nil {
				// Return episode if it's in the future or aired in the last 24 hours
				timeDiff := currentTime.UnixMilli() - airTime.UnixMilli()
				dayMs := int64(24 * 60 * 60 * 1000)

				if airTime.After(currentTime) || timeDiff <= dayMs {
					animeID := ""
					if episode.AnimeID != nil {
						animeID = *episode.AnimeID
					}
					return &NextEpisodeResult{
						Episode: Episode{
							ID:            episode.ID,
							AnimeID:       animeID,
							EpisodeNumber: episode.EpisodeNumber,
							TitleEn:       episode.TitleEn,
							TitleJp:       episode.TitleJp,
							AirDate:       episode.AirDate,
							Synopsis:      episode.Synopsis,
						},
						AirTime: *airTime,
					}
				}
			}
		}
	}

	return nil
}

// GetAirTimeDisplay gets air time display configuration for AnimeCard component
func GetAirTimeDisplay(airDate *time.Time, broadcast *string, duration *string, currentTime time.Time) *AirTimeDisplayInfo {
	if airDate == nil {
		return nil
	}

	durationMinutes := parseDurationToMinutes(duration)

	// Check if currently airing first
	if isCurrentlyAiring(airDate, broadcast, durationMinutes, currentTime) {
		countdown := calculateCountdown(airDate, broadcast, durationMinutes, currentTime)
		text := "Airing"
		if countdown != "" && countdown != "AIRING NOW" {
			text = fmt.Sprintf("Airing (%s)", countdown)
		}
		return &AirTimeDisplayInfo{
			Show:    true,
			Text:    text,
			Variant: "airing",
		}
	}

	airTime := ParseAirTime(airDate, broadcast)
	if airTime == nil {
		return nil
	}

	// Check if airing today
	timeDiff := airTime.UnixMilli() - currentTime.UnixMilli()
	dayMs := int64(24 * 60 * 60 * 1000)
	isAiringToday := timeDiff > 0 && timeDiff <= dayMs

	if isAiringToday {
		countdown := calculateCountdown(airDate, broadcast, durationMinutes, currentTime)
		if countdown != "" {
			isJustAired := countdown == "JUST AIRED"
			text := "Just aired"
			variant := "aired"

			if !isJustAired {
				if strings.Contains(countdown, "AIRING NOW") {
					text = "Airing now"
					variant = "countdown"
				} else {
					text = fmt.Sprintf("Airing in %s", countdown)
					variant = "countdown"
				}
			}

			return &AirTimeDisplayInfo{
				Show:    true,
				Text:    text,
				Variant: variant,
			}
		}
	}

	if hasAlreadyAired(airDate, broadcast, durationMinutes, currentTime) {
		return &AirTimeDisplayInfo{
			Show:    true,
			Text:    "Recently aired",
			Variant: "aired",
		}
	}

	// For scheduled episodes, show shorter format for cards
	shortDateTime := airTime.Format("Mon at 3:04 PM")

	// Check if episode is within 24 hours to decide whether to show "Airing" prefix
	isWithin24Hours := timeDiff <= dayMs

	text := shortDateTime
	if isWithin24Hours {
		text = fmt.Sprintf("Airing %s", shortDateTime)
	}

	return &AirTimeDisplayInfo{
		Show:    true,
		Text:    text,
		Variant: "scheduled",
	}
}

// ProcessCurrentlyAiring processes currently airing data like the frontend
func ProcessCurrentlyAiring(animes []*model.Anime, limit int, currentTime time.Time) []*model.Anime {
	if len(animes) == 0 {
		return []*model.Anime{}
	}

	type ProcessedAnime struct {
		Anime           *model.Anime
		NextEpisodeDate time.Time
		NextEpisode     *Episode
		AirTimeDisplay  *AirTimeDisplayInfo
	}

	var processedAnime []ProcessedAnime

	// Process each anime and determine its next episode
	for _, anime := range animes {
		if anime.Episodes == nil || len(anime.Episodes) == 0 {
			continue
		}

		// Find the next episode
		nextEpisodeResult := findNextEpisode(anime.Episodes, anime.Broadcast, currentTime)
		if nextEpisodeResult == nil {
			continue
		}

		// Generate air time display info
		airTimeInfo := GetAirTimeDisplay(nextEpisodeResult.Episode.AirDate, anime.Broadcast, anime.Duration, currentTime)
		if airTimeInfo == nil {
			continue
		}

		// Create processed entry
		processedEntry := ProcessedAnime{
			Anime:           anime,
			NextEpisodeDate: nextEpisodeResult.AirTime,
			NextEpisode:     &nextEpisodeResult.Episode,
			AirTimeDisplay:  airTimeInfo,
		}

		processedAnime = append(processedAnime, processedEntry)
	}

	// Filter and sort anime: only recently aired (last 30 minutes) or future episodes
	thirtyMinutesAgo := currentTime.Add(-30 * time.Minute)

	var recentlyAired []ProcessedAnime
	var futureEpisodes []ProcessedAnime

	for _, anime := range processedAnime {
		airTime := anime.NextEpisodeDate
		if airTime.Before(currentTime) && airTime.After(thirtyMinutesAgo) {
			recentlyAired = append(recentlyAired, anime)
		} else if airTime.After(currentTime) {
			futureEpisodes = append(futureEpisodes, anime)
		}
	}

	// Sort recently aired (most recent first)
	sort.Slice(recentlyAired, func(i, j int) bool {
		return recentlyAired[i].NextEpisodeDate.After(recentlyAired[j].NextEpisodeDate)
	})

	// Sort future episodes (earliest first)
	sort.Slice(futureEpisodes, func(i, j int) bool {
		return futureEpisodes[i].NextEpisodeDate.Before(futureEpisodes[j].NextEpisodeDate)
	})

	// Limit recently aired to 2 shows
	if len(recentlyAired) > 2 {
		recentlyAired = recentlyAired[:2]
	}

	// Combine: recently aired first, then future episodes
	var result []ProcessedAnime
	result = append(result, recentlyAired...)
	result = append(result, futureEpisodes...)

	// Apply limit
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	// Convert back to []*model.Anime
	var finalResult []*model.Anime
	for _, processedItem := range result {
		// Calculate the air time with timezone conversion
		var airTime *time.Time
		if processedItem.NextEpisode.AirDate != nil {
			// Use the broadcast info from the anime to calculate the proper air time
			calculatedAirTime := ParseAirTime(processedItem.NextEpisode.AirDate, processedItem.Anime.Broadcast)
			airTime = calculatedAirTime
		}

		// Update the anime with the next episode information
		processedItem.Anime.NextEpisode = &model.Episode{
			ID:            processedItem.NextEpisode.ID,
			AnimeID:       &processedItem.NextEpisode.AnimeID,
			EpisodeNumber: processedItem.NextEpisode.EpisodeNumber,
			TitleEn:       processedItem.NextEpisode.TitleEn,
			TitleJp:       processedItem.NextEpisode.TitleJp,
			AirDate:       processedItem.NextEpisode.AirDate,
			AirTime:       airTime,
			Synopsis:      processedItem.NextEpisode.Synopsis,
			CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		}

		finalResult = append(finalResult, processedItem.Anime)
	}

	return finalResult
}