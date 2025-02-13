package router

import (
	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func setupPagesRouter(apiCfg handlers.APIConfig) *chi.Mux {
	pagesRouter := chi.NewRouter()
	pagesRouter.Get("/login", handlers.LoginPageHandler)
	pagesRouter.Get("/register", handlers.RegisterPageHandler)
	pagesRouter.Get("/dashboard", apiCfg.MiddlewareAuth(handlers.DashboardPageHandler))

	return pagesRouter
}
