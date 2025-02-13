package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/eddietindame/gorssagg/internal/rss"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in env")
	}
	dbUrl := os.Getenv("DB_URL")
	if port == "" {
		log.Fatal("DB_URL not found in env")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}

	db := database.New(conn)
	apiCfg := handlers.APIConfig{
		DB: db,
	}

	go rss.StartScraping(db, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://", "http://"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server running on port %v\n", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
