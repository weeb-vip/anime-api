package anime_character_staff_link

import (
	"context"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character_staff_link"
)

type AnimeCharacterStaffLinkImpl interface {
	FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*anime_character_staff_link.AnimeCharacterWithStaff, error)
}

type AnimeCharacterStaffLinkService struct {
	Repository anime_character_staff_link.AnimeCharacterStaffLinkRepositoryImpl
}

func NewAnimeCharacterStaffLinkService(repository anime_character_staff_link.AnimeCharacterStaffLinkRepositoryImpl) AnimeCharacterStaffLinkImpl {
	return &AnimeCharacterStaffLinkService{
		Repository: repository,
	}
}

func (a *AnimeCharacterStaffLinkService) FindAnimeCharacterAndStaffByAnimeId(ctx context.Context, animeId string) ([]*anime_character_staff_link.AnimeCharacterWithStaff, error) {
	animeCharacters, err := a.Repository.FindAnimeCharacterAndStaffByAnimeId(ctx, animeId)
	if err != nil {
		return nil, err
	}
	return animeCharacters, nil
}
