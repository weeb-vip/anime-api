package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/graph"
	"github.com/weeb-vip/anime-api/graph/generated"
	"github.com/weeb-vip/anime-api/internal/db"
	anime2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character"
	"github.com/weeb-vip/anime-api/internal/db/repositories/anime_character_staff_link"
	anime3 "github.com/weeb-vip/anime-api/internal/db/repositories/anime_episode"
	"github.com/weeb-vip/anime-api/internal/directives"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	anime_character2 "github.com/weeb-vip/anime-api/internal/services/anime_character"
	anime_character_staff_link2 "github.com/weeb-vip/anime-api/internal/services/anime_character_staff_link"
	"github.com/weeb-vip/anime-api/internal/services/episodes"
	"net/http"
)

func BuildRootHandler(conf config.Config) http.Handler {
	database := db.NewDatabase(conf.DBConfig)
	animeRepository := anime2.NewAnimeRepository(database)
	episodeRepository := anime3.NewAnimeEpisodeRepository(database)
	animeService := anime.NewAnimeService(animeRepository)
	animeEpisodeService := episodes.NewAnimeEpisodeService(episodeRepository)
	animeCharacterRepository := anime_character.NewAnimeCharacterRepository(database)
	animeCharacterService := anime_character2.NewAnimeCharacterService(animeCharacterRepository)
	animeCharacterWithStaffLinkRepository := anime_character_staff_link.NewAnimeCharacterStaffLinkRepository(database)
	animeCharacterWithStaffLinkService := anime_character_staff_link2.NewAnimeCharacterStaffLinkService(animeCharacterWithStaffLinkRepository)
	resolvers := &graph.Resolver{
		Config:                             conf,
		AnimeService:                       animeService,
		AnimeEpisodeService:                animeEpisodeService,
		AnimeCharacterService:              animeCharacterService,
		AnimeCharacterWithStaffLinkService: animeCharacterWithStaffLinkService,
	}

	cfg := generated.Config{Resolvers: resolvers, Directives: directives.GetDirectives()}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	return srv
}
