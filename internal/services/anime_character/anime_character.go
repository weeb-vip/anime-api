package anime_character

import (
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
)

type AnimeCharacterServiceImpl interface {
}

type AnimeCharacterService struct {
	Repository anime_character.AnimeCharacterRepositoryImpl
}

func NewAnimeCharacterService(repository anime_character.AnimeCharacterRepositoryImpl) AnimeCharacterServiceImpl {
	return &AnimeCharacterService{
		Repository: repository,
	}
}
