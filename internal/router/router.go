package router

import (
	"net/http"

	templMiddleware "github.com/axzilla/templui/middleware"
	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
)

func SetupRouter() *chi.Mux {
	apiCfg := handlers.APIConfig{
		DB: database.GetQueries(),
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(
		csrf.Protect(
			[]byte(config.CSRF_KEY),
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

	if config.Environment == "production" {
		cspConfig := templMiddleware.CSPConfig{
			ScriptSrc: []string{"cdn.jsdelivr.net"},
		}
		router.Use(templMiddleware.WithCSP(cspConfig))
	}

	fs := http.FileServer(http.Dir("public"))
	router.Handle("/public/*", http.StripPrefix("/public/", fs))
	router.Mount("/v1", setupV1Router(apiCfg))
	router.Mount("/", setupPagesRouter(apiCfg))

	return router
}
