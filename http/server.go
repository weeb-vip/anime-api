package http

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/http/handlers"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"log"
	"net/http"
)

func SetupServer(cfg config.Config) *muxtrace.Router {

	router := muxtrace.NewRouter()

	router.Handle("/ui/playground", playground.Handler("GraphQL playground", "/graphql")).Methods("GET")
	router.Handle("/graphql", handlers.BuildRootHandler(cfg)).Methods("POST")
	router.Handle("/healthcheck", handlers.HealthCheckHandler()).Methods("GET")

	return router
}

func StartServer() error {
	cfg := config.LoadConfigOrPanic()
	router := SetupServer(cfg)
	tracer.Start()
	tracer.WithService(cfg.AppConfig.APPName)
	tracer.WithServiceVersion(cfg.AppConfig.Version)
	defer tracer.Stop()

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.AppConfig.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppConfig.Port), router)
}
