package anime_character_staff_link

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_staff"
)

type AnimeCharacterWithStaff struct {
	// Foreign keys for relationships
	CharacterID string `gorm:"type:varchar(36);not null"` // Foreign key for AnimeCharacter
	StaffID     string `gorm:"type:varchar(36);not null"` // Foreign key for AnimeStaff

	// The actual associated structs
	AnimeCharacter anime_character.AnimeCharacter `json:"anime_character" gorm:"foreignkey:CharacterID;references:ID"` // Foreign key to AnimeCharacter
	AnimeStaff     anime_staff.AnimeStaff         `json:"anime_staff" gorm:"foreignkey:StaffID;references:ID"`         // Foreign key to AnimeStaff
}

func (AnimeCharacterWithStaff) TableName() string {
	table := AnimeCharacterStaffLink{}
	return table.TableName()
}

type AnimeCharacterStaffLinkRepositoryImpl interface {
	FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*AnimeCharacterWithStaff, error)
}

type AnimeCharacterStaffLinkRepository struct {
	db *db.DB
}

func NewAnimeCharacterStaffLinkRepository(db *db.DB) AnimeCharacterStaffLinkRepositoryImpl {
	return &AnimeCharacterStaffLinkRepository{db: db}
}

func (a *AnimeCharacterStaffLinkRepository) FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*AnimeCharacterWithStaff, error) {
	var results []*AnimeCharacterWithStaff

	err := a.db.DB.
		Preload("AnimeCharacter").
		Preload("AnimeStaff").
		Joins("JOIN anime_character ON anime_character.id = anime_character_staff_link.character_id").
		Where("anime_character.anime_id = ?", animeId).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
