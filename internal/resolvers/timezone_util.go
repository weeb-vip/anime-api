package resolvers

import (
	"time"
)

// getJapanTimezone returns the Japan timezone (JST)
func getJapanTimezone() *time.Location {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// Fallback to UTC+9 if Asia/Tokyo is not available
		jst = time.FixedZone("JST", 9*60*60)
	}
	return jst
}

// convertToJapanTime converts a time.Time to Japan timezone if it's not nil
// This ensures that dates stored as local Japan time are properly displayed with JST context
func convertToJapanTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	
	jst := getJapanTimezone()
	
	// If the time doesn't have timezone info (which is likely the case for our stored dates),
	// we assume it's already in Japan local time and just assign the timezone
	japanTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), jst)
	return &japanTime
}