package anime_character

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
)

type AnimeCharacterServiceImpl interface {
	FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*anime_character.AnimeCharacterWithStaff, error)
}

type AnimeCharacterService struct {
	Repository anime_character.AnimeCharacterRepositoryImpl
}

func NewAnimeCharacterService(repository anime_character.AnimeCharacterRepositoryImpl) AnimeCharacterServiceImpl {
	return &AnimeCharacterService{
		Repository: repository,
	}
}

func (a *AnimeCharacterService) FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*anime_character.AnimeCharacterWithStaff, error) {
	animeCharacters, err := a.Repository.FindAnimeCharacterAndStaffByAnimeId(ctx, animeId)
	if err != nil {
		return nil, err
	}
	return animeCharacters, nil
}
