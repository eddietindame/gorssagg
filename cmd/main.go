package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/router"
	"github.com/eddietindame/gorssagg/internal/rss"
	"github.com/eddietindame/gorssagg/internal/store"
	_ "github.com/lib/pq"
)

func init() {
	config.InitEnv()
}

func main() {
	store.InitRedis()
	defer store.RedisClient.Close()
	store.InitSessionStore()
	defer store.Store.Close()

	db := database.GetQueries()

	go rss.StartScraping(db, 10, time.Minute)

	router := router.SetupRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + config.PORT,
	}

	fmt.Println("Server running on port", config.PORT)
	log.Fatal(srv.ListenAndServe())
}
