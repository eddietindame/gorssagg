package router

import (
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter() *chi.Mux {
	db := database.GetQueries()
	apiCfg := handlers.APIConfig{
		DB: db,
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://", "http://"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Mount("/v1", setupV1Router(apiCfg))
	router.Mount("/", setupPagesRouter(apiCfg))

	return router
}
