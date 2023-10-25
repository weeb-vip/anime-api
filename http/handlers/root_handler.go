package handlers

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/weeb-vip/anime-api/graph/generated"
	"github.com/weeb-vip/anime-api/internal/db"
	anime2 "github.com/weeb-vip/anime-api/internal/db/repositories/anime"
	"github.com/weeb-vip/anime-api/internal/services/anime"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/graph"
	"github.com/weeb-vip/anime-api/internal/directives"
	"net/http"
)

func BuildRootHandler(conf config.Config) http.Handler {
	database := db.NewDatabase(conf.DBConfig)
	animeRepository := anime2.NewAnimeRepository(database)
	animeService := anime.NewAnimeService(animeRepository)
	resolvers := &graph.Resolver{
		Config:       conf,
		AnimeService: animeService,
	}

	cfg := generated.Config{Resolvers: resolvers, Directives: directives.GetDirectives()}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(cfg))

	return srv
}

// middleware to pass span in context
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := tracer.StartSpanFromContext(r.Context(), "http.request")
		defer span.Finish()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
