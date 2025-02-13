package router

import (
	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func setupV1Router(apiCfg handlers.APIConfig) *chi.Mux {
	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handlers.Readiness)
	v1Router.Get("/err", handlers.Err)

	v1Router.Get("/users", apiCfg.MiddlewareAuth((apiCfg.GetUser)))
	v1Router.Post("/users", apiCfg.CreateUser)

	v1Router.Get("/feeds", apiCfg.GetFeeds)
	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.CreateFeed))

	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.GetPostsForUser))

	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.GetFeedFollows))
	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.CreateFeedFollow))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.DeleteFeedFollow))

	return v1Router
}
