package anime_fanart

import "time"

type Fanart struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	AnimeID   string    `gorm:"column:anime_id" json:"anime_id"`
	ImageURL  string    `gorm:"column:image_url" json:"image_url"`
	SourceURL *string   `gorm:"column:source_url" json:"source_url"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Fanart) TableName() string {
	return "anime_fanart"
}
