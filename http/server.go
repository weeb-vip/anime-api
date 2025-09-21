package http

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/http/handlers"
	"github.com/weeb-vip/anime-api/internal/logger"
	"github.com/weeb-vip/anime-api/metrics"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"net/http"
)

func SetupServer(cfg config.Config) *muxtrace.Router {

	router := muxtrace.NewRouter()
	router.Handle("/ui/playground", playground.Handler("GraphQL playground", "/graphql")).Methods("GET")
	router.Handle("/graphql", handlers.BuildRootHandler(cfg)).Methods("POST")
	router.Handle("/healthcheck", handlers.HealthCheckHandler()).Methods("GET")
	router.Handle("/metrics", metrics.NewPrometheusInstance().Handler()).Methods("GET")

	return router
}

func SetupServerWithContext(ctx context.Context, cfg config.Config) *muxtrace.Router {

	router := muxtrace.NewRouter(muxtrace.WithServiceName(cfg.AppConfig.APPName))
	router.Handle("/ui/playground", playground.Handler("GraphQL playground", "/graphql")).Methods("GET")
	router.Handle("/graphql", handlers.BuildRootHandlerWithContext(ctx, cfg)).Methods("POST")
	router.Handle("/healthcheck", handlers.HealthCheckHandler()).Methods("GET")
	router.Handle("/metrics", metrics.NewPrometheusInstance().Handler()).Methods("GET")

	return router
}

func StartServer() error {
	cfg := config.LoadConfigOrPanic()
	router := SetupServer(cfg)

	logger.Get().Info().
		Int("port", cfg.AppConfig.Port).
		Str("playground_url", fmt.Sprintf("http://localhost:%d/", cfg.AppConfig.Port)).
		Msg("Starting GraphQL server")

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppConfig.Port), router)
}

func StartServerWithContext(ctx context.Context) error {
	cfg := config.LoadConfigOrPanic()
	router := SetupServerWithContext(ctx, cfg)

	logger.FromCtx(ctx).Info().
		Int("port", cfg.AppConfig.Port).
		Str("playground_url", fmt.Sprintf("http://localhost:%d/", cfg.AppConfig.Port)).
		Msg("Starting GraphQL server")

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppConfig.Port), router)
}
