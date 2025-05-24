package anime_character

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db"
)

type AnimeCharacterRepositoryImpl interface {
	FindAnimearacterById(ctx context.Context, id string) (*AnimeCharacter, error)
}

type AnimeCharacterRepository struct {
	db *db.DB
}

func NewAnimeCharacterRepository(db *db.DB) AnimeCharacterRepositoryImpl {
	return &AnimeCharacterRepository{db: db}
}

func (a *AnimeCharacterRepository) FindAnimearacterById(ctx context.Context, id string) (*AnimeCharacter, error) {
	var animeCharacter AnimeCharacter
	err := a.db.DB.Where("id = ?", id).First(&animeCharacter).Error
	if err != nil {
		return nil, err
	}
	return &animeCharacter, nil
}
