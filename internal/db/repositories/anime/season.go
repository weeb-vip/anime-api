package anime

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Season string

var seasonPattern = regexp.MustCompile(`^(SPRING|SUMMER|FALL|WINTER)_(\d{4})$`)

func (s Season) String() string {
	return string(s)
}

func (s *Season) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	if str, ok := value.(string); ok {
		*s = Season(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into Season", value)
}

func (s Season) Value() (driver.Value, error) {
	return string(s), nil
}

func ParseSeason(s string) (Season, error) {
	if !seasonPattern.MatchString(s) {
		return "", fmt.Errorf("invalid season year format: %s (expected format: SPRING_2024)", s)
	}
	return Season(s), nil
}

func (s Season) IsValid() bool {
	return seasonPattern.MatchString(string(s))
}

func (s Season) GetSeason() string {
	matches := seasonPattern.FindStringSubmatch(string(s))
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func (s Season) GetYear() int {
	matches := seasonPattern.FindStringSubmatch(string(s))
	if len(matches) >= 3 {
		year, _ := strconv.Atoi(matches[2])
		return year
	}
	return 0
}

func CreateSeason(season string, year int) Season {
	return Season(fmt.Sprintf("%s_%d", strings.ToUpper(season), year))
}
