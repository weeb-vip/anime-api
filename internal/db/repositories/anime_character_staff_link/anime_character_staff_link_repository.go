package anime_character_staff_link

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_staff"
)

type AnimeCharacterWithStaff struct {
	anime_character.AnimeCharacter
	VoiceActors []anime_staff.AnimeStaff `json:"voice_actors"`
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
	type joinResult struct {
		CharacterID    string
		CharacterName  string
		CharacterRole  string
		CharacterImage string
		StaffID        string
		GivenName      string
		FamilyName     string
		StaffImage     string
	}

	var rows []joinResult

	err := a.db.DB.
		Table("anime_character_staff_link").
		Select(`
			anime_character.id as character_id,
			anime_character.name as character_name,
			anime_character.role as character_role,
			anime_character.image as character_image,
			anime_staff.id as staff_id,
			anime_staff.given_name,
			anime_staff.family_name,
			anime_staff.image as staff_image
		`).
		Joins("JOIN anime_character ON anime_character.id = anime_character_staff_link.character_id").
		Joins("JOIN anime_staff ON anime_staff.id = anime_character_staff_link.staff_id").
		Where("anime_character.anime_id = ?", animeId).
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	// Group by character ID
	characterMap := make(map[string]*AnimeCharacterWithStaff)
	for _, row := range rows {
		if _, exists := characterMap[row.CharacterID]; !exists {
			characterMap[row.CharacterID] = &AnimeCharacterWithStaff{
				AnimeCharacter: anime_character.AnimeCharacter{
					ID:    row.CharacterID,
					Name:  row.CharacterName,
					Role:  row.CharacterRole,
					Image: row.CharacterImage,
				},
				VoiceActors: []anime_staff.AnimeStaff{},
			}
		}

		characterMap[row.CharacterID].VoiceActors = append(characterMap[row.CharacterID].VoiceActors, anime_staff.AnimeStaff{
			ID:         row.StaffID,
			GivenName:  row.GivenName,
			FamilyName: row.FamilyName,
			Image:      row.StaffImage,
		})
	}

	var result []*AnimeCharacterWithStaff
	for _, char := range characterMap {
		result = append(result, char)
	}

	return result, nil
}
