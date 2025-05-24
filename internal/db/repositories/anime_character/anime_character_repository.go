package anime_character

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_staff"
)

type AnimeCharacterWithStaff struct {
	AnimeCharacter *AnimeCharacter         `json:"anime_character"`
	AnimeStaff     *anime_staff.AnimeStaff `json:"anime_staff"`
}

type AnimeCharacterRepositoryImpl interface {
	FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*AnimeCharacterWithStaff, error)
	FindAnimearacterById(ctx context.Context, id string) (*AnimeCharacter, error)
}

type AnimeCharacterRepository struct {
	db *db.DB
}

func NewAnimeCharacterRepository(db *db.DB) AnimeCharacterRepositoryImpl {
	return &AnimeCharacterRepository{db: db}
}

func (a *AnimeCharacterRepository) FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*AnimeCharacterWithStaff, error) {
	var animeCharactersWithStaff []*AnimeCharacterWithStaff
	// use anime_character_staff_link table to join anime_character and anime_staff
	err := a.db.DB.Table("anime_character").
		Select("anime_character.*, anime_staff.given_name, anime_staff.family_name, anime_staff.birth_place, anime_staff.blood_type, anime_staff.hobbies").
		Joins("JOIN anime_character_staff_link ON anime_character.id = anime_character_staff_link.character_id").
		Joins("JOIN anime_staff ON anime_character_staff_link.staff_id = anime_staff.id").
		Where("anime_character.anime_id = ?", animeId).
		Scan(&animeCharactersWithStaff).Error
	if err != nil {
		return nil, err
	}
	return animeCharacters, nil
}

func (a *AnimeCharacterRepository) FindAnimearacterById(ctx context.Context, id string) (*AnimeCharacter, error) {
	var animeCharacter AnimeCharacter
	err := a.db.DB.Where("id = ?", id).First(&animeCharacter).Error
	if err != nil {
		return nil, err
	}
	return &animeCharacter, nil
}
