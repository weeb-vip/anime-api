package http

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/http/handlers"
	"github.com/weeb-vip/anime-api/metrics"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"log"
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

func StartServer() error {
	cfg := config.LoadConfigOrPanic()
	router := SetupServer(cfg)

	log.Printf("connect to http://localhost:%d/ for GraphQL playground", cfg.AppConfig.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppConfig.Port), router)
}
