package store

import (
	"log"
	"net/http"

	"github.com/boj/redistore"
	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/gorilla/sessions"
)

const Authenticated string = "authenticated"
const UserID string = "userId"
const Email string = "email"
const Username string = "username"

var Store *redistore.RediStore

func InitSessionStore() {
	var err error
	Store, err = redistore.NewRediStore(
		10,
		"tcp",
		config.REDIS_HOST,
		config.REDIS_USERNAME,
		config.REDIS_PASSWORD,
		[]byte(config.SESSION_KEY),
	)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	Store.Options = &sessions.Options{
		HttpOnly: true,                               // Prevent JavaScript from accessing the cookie
		Secure:   config.Environment == "production", // Send only over HTTPS
		SameSite: http.SameSiteStrictMode,            // Prevent CSRF attacks
		MaxAge:   3600,                               // Session expires in 1 hour
	}
}
