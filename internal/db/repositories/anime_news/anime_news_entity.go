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
	// ISO 639-1 code of the SOURCE article. The summary is always English, so this is
	// what tells a reader the link itself is Japanese before they follow it.
	Language *string `gorm:"column:language" json:"language"`
	// JSON array of {kind,title,url} for media the article points at (PV, official
	// site, announcement post). The column is `reference_links` because REFERENCES is
	// reserved in MySQL and would need backticking everywhere it appeared.
	ReferenceLinks *string   `gorm:"column:reference_links" json:"reference_links"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (AnimeNews) TableName() string {
	return "anime_news"
}
