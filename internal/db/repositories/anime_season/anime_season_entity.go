package anime_season

import (
	"time"
)

type AnimeSeasonStatus string

const (
	StatusUnknown   AnimeSeasonStatus = "unknown"
	StatusConfirmed AnimeSeasonStatus = "confirmed"
	StatusAnnounced AnimeSeasonStatus = "announced"
	StatusCancelled AnimeSeasonStatus = "cancelled"
)

type AnimeSeason struct {
	ID           string            `gorm:"column:id;type:varchar(36);primaryKey;default:(UUID())" json:"id"`
	Season       string            `gorm:"column:season;not null" json:"season"`
	Status       AnimeSeasonStatus `gorm:"column:status;type:enum('unknown','confirmed','announced','cancelled');default:unknown;not null" json:"status"`
	EpisodeCount *int              `gorm:"column:episode_count;null" json:"episode_count"`
	Notes        *string           `gorm:"column:notes;type:text;null" json:"notes"`
	CreatedAt    time.Time         `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt    time.Time         `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;not null" json:"updated_at"`
	AnimeID      *string           `gorm:"column:anime_id;type:varchar(36);null" json:"anime_id"`
}

// TableName sets the table name
func (AnimeSeason) TableName() string {
	return "anime_seasons"
}