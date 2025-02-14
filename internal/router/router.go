package router

import (
	"log"
	"os"

	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
)

func SetupRouter() *chi.Mux {
	csrfKey := os.Getenv("CSRF_KEY")
	if csrfKey == "" {
		log.Fatal("CSRF_KEY not found in env")
	}

	db := database.GetQueries()
	apiCfg := handlers.APIConfig{
		DB: db,
	}

	router := chi.NewRouter()
	router.Use(
		csrf.Protect(
			[]byte(csrfKey),
			csrf.Secure(config.Environment == "production"),
			csrf.HttpOnly(true),
			csrf.SameSite(csrf.SameSiteStrictMode),
			csrf.Path("/"),
		),
	)
	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://", "http://"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}),
	)

	router.Mount("/v1", setupV1Router(apiCfg))
	router.Mount("/", setupPagesRouter(apiCfg))

	return router
}
