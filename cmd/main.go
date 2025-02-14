package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers"
	"github.com/eddietindame/gorssagg/internal/router"
	"github.com/eddietindame/gorssagg/internal/rss"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in env")
	}
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("SESSION_KEY not found in env")
	}
	handlers.Store = sessions.NewCookieStore([]byte(sessionKey))
	handlers.Store.Options = &sessions.Options{Path: "/",
		HttpOnly: true,                               // Prevent JavaScript from accessing the cookie
		Secure:   config.Environment == "production", // Send only over HTTPS
		SameSite: http.SameSiteStrictMode,            // Prevent CSRF attacks
		MaxAge:   3600,                               // Session expires in 1 hour
	}

	db := database.GetQueries()

	go rss.StartScraping(db, 10, time.Minute)

	router := router.SetupRouter()
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Println("Server running on port", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
