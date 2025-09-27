package anime

import (
	"testing"
	"time"
)

func TestDateConversions(t *testing.T) {
	// Load Tokyo timezone for testing
	tzTokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal("Failed to load Tokyo timezone:", err)
	}

	tests := []struct {
		name           string
		inputUTC       string
		expectedJST    string
		description    string
	}{
		{
			name:        "UTC midnight to JST",
			inputUTC:    "2025-09-26T00:00:00Z",
			expectedJST: "2025-09-26 09:00:00 +0900 JST",
			description: "UTC midnight should be 9AM JST (UTC+9)",
		},
		{
			name:        "Show airing at 00:30 JST on Sept 28",
			inputUTC:    "2025-09-27T15:30:00Z",
			expectedJST: "2025-09-28 00:30:00 +0900 JST",
			description: "15:30 UTC on Sept 27 should be 00:30 JST on Sept 28",
		},
		{
			name:        "UTC 08:00 to JST",
			inputUTC:    "2025-09-27T08:00:00Z",
			expectedJST: "2025-09-27 17:00:00 +0900 JST",
			description: "8AM UTC should be 5PM JST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input UTC time
			utcTime, err := time.Parse(time.RFC3339, tt.inputUTC)
			if err != nil {
				t.Fatalf("Failed to parse UTC time %s: %v", tt.inputUTC, err)
			}

			// Convert to JST
			jstTime := utcTime.In(tzTokyo)

			// Check if the converted time matches expected
			if jstTime.String() != tt.expectedJST {
				t.Errorf("Time conversion mismatch\n"+
					"Description: %s\n"+
					"Input UTC:   %s\n"+
					"Expected JST: %s\n"+
					"Got JST:     %s",
					tt.description, tt.inputUTC, tt.expectedJST, jstTime.String())
			}
		})
	}
}

func TestDateRangeWithDays(t *testing.T) {
	tzTokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal("Failed to load Tokyo timezone:", err)
	}

	tests := []struct {
		name         string
		startUTC     string
		days         int
		expectedEnd  string
		shouldInclude []string // Times that should be included in the range
		shouldExclude []string // Times that should NOT be included
	}{
		{
			name:        "7 days from Sept 26 midnight UTC",
			startUTC:    "2025-09-26T00:00:00Z",
			days:        7,
			expectedEnd: "2025-10-03 09:00:00 +0900 JST",
			shouldInclude: []string{
				"2025-09-26T00:00:00Z", // Start of range
				"2025-09-27T15:30:00Z", // Show at 00:30 JST Sept 28
				"2025-09-27T08:00:00Z", // Show at 17:00 JST Sept 27
				"2025-10-02T23:59:59Z", // Just before end
			},
			shouldExclude: []string{
				"2025-09-25T23:59:59Z", // Just before start
				"2025-10-03T00:00:01Z", // Just after end
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime, err := time.Parse(time.RFC3339, tt.startUTC)
			if err != nil {
				t.Fatalf("Failed to parse start time: %v", err)
			}

			// Convert to JST as the code does
			startJST := startTime.In(tzTokyo)
			endJST := startJST.AddDate(0, 0, tt.days)

			if endJST.String() != tt.expectedEnd {
				t.Errorf("End date mismatch\nExpected: %s\nGot: %s",
					tt.expectedEnd, endJST.String())
			}

			// Check included times
			for _, timeStr := range tt.shouldInclude {
				checkTime, _ := time.Parse(time.RFC3339, timeStr)
				checkJST := checkTime.In(tzTokyo)

				if !isInRange(checkJST, startJST, endJST) {
					t.Errorf("Time should be included but wasn't:\n"+
						"Time: %s (UTC) -> %s (JST)\n"+
						"Range: %s to %s (JST)",
						timeStr, checkJST, startJST, endJST)
				}
			}

			// Check excluded times
			for _, timeStr := range tt.shouldExclude {
				checkTime, _ := time.Parse(time.RFC3339, timeStr)
				checkJST := checkTime.In(tzTokyo)

				if isInRange(checkJST, startJST, endJST) {
					t.Errorf("Time should be excluded but wasn't:\n"+
						"Time: %s (UTC) -> %s (JST)\n"+
						"Range: %s to %s (JST)",
						timeStr, checkJST, startJST, endJST)
				}
			}
		})
	}
}

func TestStartOfDayIn(t *testing.T) {
	tzTokyo, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal("Failed to load Tokyo timezone:", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "UTC midnight",
			input:    "2025-09-26T00:00:00Z",
			expected: "2025-09-26 00:00:00 +0900 JST",
		},
		{
			name:     "UTC noon",
			input:    "2025-09-26T12:00:00Z",
			expected: "2025-09-26 00:00:00 +0900 JST",
		},
		{
			name:     "UTC late evening (next day in JST)",
			input:    "2025-09-26T16:00:00Z",
			expected: "2025-09-27 00:00:00 +0900 JST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputTime, _ := time.Parse(time.RFC3339, tt.input)
			result := startOfDayIn(inputTime, tzTokyo)

			if result.String() != tt.expected {
				t.Errorf("startOfDayIn mismatch\n"+
					"Input:    %s\n"+
					"Expected: %s\n"+
					"Got:      %s",
					tt.input, tt.expected, result.String())
			}
		})
	}
}

// Helper function to check if a time is within a range (inclusive)
func isInRange(t, start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Before(end) || t.Equal(end))
}