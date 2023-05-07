package graph

import (
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/internal/services/anime"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Config       config.Config
	AnimeService anime.AnimeServiceImpl
}
