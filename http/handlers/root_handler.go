package handlers

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/graph"
	"github.com/weeb-vip/anime-api/graph/generated"
	"github.com/weeb-vip/anime-api/internal/cache"
	"github.com/weeb-vip/anime-api/internal/db"
	anime2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character_staff_link"
	anime3 "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_season"
	"github.com/weeb-vip/anime-api/internal/directives"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	anime_character2 "github.com/weeb-vip/anime-api/internal/services/anime_character"
	anime_character_staff_link2 "github.com/weeb-vip/anime-api/internal/services/anime_character_staff_link"
	anime_season_service "github.com/weeb-vip/anime-api/internal/services/anime_season"
	"github.com/weeb-vip/anime-api/internal/services/episodes"
	"github.com/weeb-vip/anime-api/http/middleware"
	"github.com/weeb-vip/anime-api/internal/logger"
	"net/http"
)

func BuildRootHandler(conf config.Config) http.Handler {
	database := db.NewDatabase(conf.DBConfig)

	// Initialize cache if enabled
	log := logger.FromCtx(context.Background())
	log.Info().
		Bool("redis_enabled", conf.RedisConfig.Enabled).
		Str("redis_host", conf.RedisConfig.Host).
		Str("redis_port", conf.RedisConfig.Port).
		Int("redis_db", conf.RedisConfig.DB).
		Msg("Cache configuration")

	cacheInstance, err := cache.NewCache(conf)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize cache, continuing without caching")
		cacheInstance = cache.NewNoOpCache()
	} else if conf.RedisConfig.Enabled {
		log.Info().Msg("Redis cache successfully initialized")
	} else {
		log.Info().Msg("Cache disabled by configuration")
	}
	cacheService := cache.NewCacheService(cacheInstance, conf.RedisConfig)

	// Initialize repositories
	var animeRepository anime2.AnimeRepositoryImpl
	var episodeRepository anime3.AnimeEpisodeRepositoryImpl

	log.Info().Bool("cache_enabled", conf.RedisConfig.Enabled).Msg("Cache configuration status")

	if conf.RedisConfig.Enabled {
		log.Info().Msg("Cache enabled, using repositories with caching")

		// Use repositories with caching when enabled
		animeRepository = anime2.NewAnimeRepositoryWithCache(database, cacheService)
		episodeRepository = anime3.NewAnimeEpisodeRepositoryWithCache(database, cacheService)
	} else {
		log.Info().Msg("Cache disabled, using direct database repositories")

		// Use direct repositories when caching is disabled
		animeRepository = anime2.NewAnimeRepository(database)
		episodeRepository = anime3.NewAnimeEpisodeRepository(database)
	}

	animeService := anime.NewAnimeService(animeRepository)
	animeEpisodeService := episodes.NewAnimeEpisodeService(episodeRepository)
	animeCharacterRepository := anime_character.NewAnimeCharacterRepository(database)
	animeCharacterService := anime_character2.NewAnimeCharacterService(animeCharacterRepository)
	animeCharacterWithStaffLinkRepository := anime_character_staff_link.NewAnimeCharacterStaffLinkRepository(database)
	animeCharacterWithStaffLinkService := anime_character_staff_link2.NewAnimeCharacterStaffLinkService(animeCharacterWithStaffLinkRepository)
	animeSeasonRepository := anime_season.NewAnimeSeasonRepository(database)
	animeSeasonService := anime_season_service.NewAnimeSeasonService(animeSeasonRepository)
	resolvers := &graph.Resolver{
		Config:                             conf,
		AnimeService:                       animeService,
		AnimeEpisodeService:                animeEpisodeService,
		AnimeCharacterService:              animeCharacterService,
		AnimeCharacterWithStaffLinkService: animeCharacterWithStaffLinkService,
		AnimeSeasonService:                 animeSeasonService,
	}

	cfg := generated.Config{Resolvers: resolvers, Directives: directives.GetDirectives()}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	// Add GraphQL tracing extension
	srv.Use(&middleware.GraphQLTracingExtension{})

	return srv
}

func BuildRootHandlerWithContext(ctx context.Context, conf config.Config) http.Handler {
	database := db.NewDatabase(conf.DBConfig)

	// Initialize cache if enabled
	cacheInstance, err := cache.NewCache(conf)
	if err != nil {
		log := logger.FromCtx(ctx)
		log.Error().Err(err).Msg("Failed to initialize cache, continuing without caching")
		cacheInstance = cache.NewNoOpCache()
	}
	cacheService := cache.NewCacheService(cacheInstance, conf.RedisConfig)

	// Initialize repositories
	var animeRepository anime2.AnimeRepositoryImpl
	var episodeRepository anime3.AnimeEpisodeRepositoryImpl

	log := logger.FromCtx(ctx)
	log.Info().Bool("cache_enabled", conf.RedisConfig.Enabled).Msg("Cache configuration status")

	if conf.RedisConfig.Enabled {
		log.Info().Msg("Cache enabled, using repositories with caching")

		// Use repositories with caching when enabled
		animeRepository = anime2.NewAnimeRepositoryWithCache(database, cacheService)
		episodeRepository = anime3.NewAnimeEpisodeRepositoryWithCache(database, cacheService)
	} else {
		log.Info().Msg("Cache disabled, using direct database repositories")

		// Use direct repositories when caching is disabled
		animeRepository = anime2.NewAnimeRepository(database)
		episodeRepository = anime3.NewAnimeEpisodeRepository(database)
	}

	animeService := anime.NewAnimeService(animeRepository)
	animeEpisodeService := episodes.NewAnimeEpisodeService(episodeRepository)
	animeCharacterRepository := anime_character.NewAnimeCharacterRepository(database)
	animeCharacterService := anime_character2.NewAnimeCharacterService(animeCharacterRepository)
	animeCharacterWithStaffLinkRepository := anime_character_staff_link.NewAnimeCharacterStaffLinkRepository(database)
	animeCharacterWithStaffLinkService := anime_character_staff_link2.NewAnimeCharacterStaffLinkService(animeCharacterWithStaffLinkRepository)
	animeSeasonRepository := anime_season.NewAnimeSeasonRepository(database)
	animeSeasonService := anime_season_service.NewAnimeSeasonService(animeSeasonRepository)
	resolvers := &graph.Resolver{
		Config:                             conf,
		AnimeService:                       animeService,
		AnimeEpisodeService:                animeEpisodeService,
		AnimeCharacterService:              animeCharacterService,
		AnimeCharacterWithStaffLinkService: animeCharacterWithStaffLinkService,
		AnimeSeasonService:                 animeSeasonService,
		Context:                            ctx,
	}

	cfg := generated.Config{Resolvers: resolvers, Directives: directives.GetDirectives()}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	// Add GraphQL tracing extension
	srv.Use(&middleware.GraphQLTracingExtension{})

	return srv
}
