package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/weeb-vip/anime-api/graph/generated"

	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/graph"
	"github.com/weeb-vip/anime-api/internal/directives"
	"net/http"
)

func BuildRootHandler(conf config.Config) http.Handler {
	resolvers := &graph.Resolver{
		Config: conf,
	}

	cfg := generated.Config{Resolvers: resolvers, Directives: directives.GetDirectives()}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	return srv
}
