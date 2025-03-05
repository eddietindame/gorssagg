package router

import (
	"fmt"

	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func setupPagesRouter(apiCfg handlers.APIConfig) *chi.Mux {
	pagesRouter := chi.NewRouter()
	pagesRouter.Get("/login", handlers.LoginPageHandler)
	pagesRouter.Post("/login", apiCfg.LoginHandler)
	pagesRouter.Get("/logout", handlers.LogoutHandler)
	pagesRouter.Get("/register", handlers.RegisterPageHandler)
	pagesRouter.Post("/register", apiCfg.RegisterHandler)
	pagesRouter.Get("/forgot-password", handlers.ForgotPageHandler)
	pagesRouter.Post("/forgot-password", apiCfg.ForgotPasswordHandler)
	pagesRouter.Get("/reset-password", handlers.ResetPageHandler)
	pagesRouter.Post("/reset-password", apiCfg.ResetPasswordHandler)
	pagesRouter.With(handlers.MiddlewareSessionAuth).Get("/posts", apiCfg.PostsPageHandler)
	pagesRouter.With(handlers.MiddlewareSessionAuth).Get("/feeds", apiCfg.FeedsPageHandler)
	pagesRouter.With(handlers.MiddlewareSessionAuth).Post("/feeds", apiCfg.FeedHandler)
	pagesRouter.With(handlers.MiddlewareSessionAuth).Delete(fmt.Sprintf("/follows/{%s}", handlers.FeedFollowsIdPageParam), apiCfg.DeleteFeedFollowHandler)

	return pagesRouter
}
