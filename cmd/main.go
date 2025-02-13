package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/router"
	"github.com/eddietindame/gorssagg/internal/rss"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in env")
	}

	db := database.GetQueries()

	go rss.StartScraping(db, 10, time.Minute)

	router := router.SetupRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server running on port %v\n", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
