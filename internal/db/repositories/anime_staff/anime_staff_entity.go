package anime_staff

import (
	"time"
)

type AnimeStaff struct {
	ID         string    `gorm:"type:char(36);primaryKey"`
	Language   string    `gorm:"type:varchar(30);not null"`
	GivenName  string    `gorm:"type:varchar(255);not null"`
	FamilyName string    `gorm:"type:varchar(255);not null"`
	// URLSlug is the public URL segment, derived in MySQL from given_name and
	// family_name by a stored generated column (migration 000038). Read-only
	// here: naming it in a write would make MySQL reject the statement.
	URLSlug *string `gorm:"column:url_slug;->"`
	Image      string    `gorm:"type:text"`
	Birthday   string    `gorm:"type:varchar(255)"`
	BirthPlace string    `gorm:"type:varchar(255)"`
	BloodType  string    `gorm:"type:varchar(255)"`
	Hobbies    string    `gorm:"type:varchar(255)"`
	Summary    string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (AnimeStaff) TableName() string {
	return "anime_staff"
}
