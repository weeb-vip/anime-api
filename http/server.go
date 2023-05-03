package http

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/weeb-vip/anime-api/config"
	"github.com/weeb-vip/anime-api/http/handlers"
	"log"
	"net/http"
)

func SetupServer(cfg config.Config) *mux.Router {

	router := mux.NewRouter()

	router.Handle("/ui/playground", playground.Handler("GraphQL playground", "/graphql")).Methods("GET")
	router.Handle("/graphql", handlers.BuildRootHandler(cfg)).Methods("POST")
	router.Handle("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	return router
}

func StartServer() error {
	cfg := config.LoadConfigOrPanic()
	router := SetupServer(cfg)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.AppConfig.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.AppConfig.Port), router)
}
