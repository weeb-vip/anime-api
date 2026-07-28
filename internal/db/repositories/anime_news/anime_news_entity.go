package anime_news

import "time"

type AnimeNews struct {
	ID            string     `gorm:"column:id;primaryKey" json:"id"`
	AnimeID       string     `gorm:"column:anime_id" json:"anime_id"`
	MalID         *int       `gorm:"column:mal_id" json:"mal_id"`
	Title         string     `gorm:"column:title" json:"title"`
	Summary       *string    `gorm:"column:summary" json:"summary"`
	Category      string     `gorm:"column:category" json:"category"`
	SourceURL     *string    `gorm:"column:source_url" json:"source_url"`
	SourceName    *string    `gorm:"column:source_name" json:"source_name"`
	PublishedDate *time.Time `gorm:"column:published_date" json:"published_date"`
	EpisodeNumber *int       `gorm:"column:episode_number" json:"episode_number"`
	TitleSlug     *string    `gorm:"column:title_slug" json:"title_slug"`
	ResearchedAt  *time.Time `gorm:"column:researched_at" json:"researched_at"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (AnimeNews) TableName() string {
	return "anime_news"
}
