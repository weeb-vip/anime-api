package graph

import (
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime_character"
	anime_character_staff_link2 "github.com/weeb-vip/anime-api/internal/services/anime_character_staff_link"
	"github.com/weeb-vip/anime-api/internal/services/episodes"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Config                             config.Config
	AnimeService                       anime.AnimeServiceImpl
	AnimeEpisodeService                episodes.AnimeEpisodeServiceImpl
	AnimeCharacterService              anime_character.AnimeCharacterServiceImpl
	AnimeCharacterWithStaffLinkService anime_character_staff_link2.AnimeCharacterStaffLinkImpl
}
