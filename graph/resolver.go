package graph

import (
	"context"
	"time"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/internal/cache"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime_character"
	anime_character_staff_link2 "github.com/weeb-vip/anime-api/internal/services/anime_character_staff_link"
	"github.com/weeb-vip/anime-api/internal/services/anime_season"
	"github.com/weeb-vip/anime-api/internal/services/episodes"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// CacheServiceInterface defines the methods needed for caching
type CacheServiceInterface interface {
	GetJSON(ctx context.Context, key string, dest interface{}) error
	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	GetKeyBuilder() *cache.CacheKeyBuilder
	GetCurrentlyAiringTTL() time.Duration
}

type Resolver struct {
	Config                             config.Config
	AnimeService                       anime.AnimeServiceImpl
	AnimeEpisodeService                episodes.AnimeEpisodeServiceImpl
	AnimeCharacterService              anime_character.AnimeCharacterServiceImpl
	AnimeCharacterWithStaffLinkService anime_character_staff_link2.AnimeCharacterStaffLinkImpl
	AnimeSeasonService                 anime_season.AnimeSeasonServiceImpl
	CacheService                       CacheServiceInterface
	Context                            context.Context
}
