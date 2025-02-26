package config

import (
	"cmp"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Set to production at build time
var Environment = "development"
var PORT, DB_URL, POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, REDIS_HOST, REDIS_USERNAME, REDIS_PASSWORD, SESSION_KEY, CSRF_KEY string

func InitEnv() {
	godotenv.Load()

	PORT = os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT not found in env")
	}

	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT = cmp.Or(os.Getenv("POSTGRES_PORT"), "5432")
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB = os.Getenv("POSTGRES_DB")

	if POSTGRES_USER != "" &&
		POSTGRES_PASSWORD != "" &&
		POSTGRES_HOST != "" &&
		POSTGRES_DB != "" {
		DB_URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			POSTGRES_USER,
			POSTGRES_PASSWORD,
			POSTGRES_HOST,
			POSTGRES_PORT,
			POSTGRES_DB,
		)
	} else {
		DB_URL = os.Getenv("DB_URL")
	}
	if DB_URL == "" {
		log.Fatal("DB_URL not found in env or could not be constructed")
	}

	REDIS_HOST = os.Getenv("REDIS_HOST")
	if REDIS_HOST == "" {
		log.Fatal("REDIS_HOST not found in env")
	}

	REDIS_USERNAME = os.Getenv("REDIS_USERNAME")
	if REDIS_USERNAME == "" && Environment == "production" {
		log.Fatal("REDIS_USERNAME not found in env")
	}

	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	if REDIS_PASSWORD == "" && Environment == "production" {
		log.Fatal("REDIS_PASSWORD not found in env")
	}

	SESSION_KEY = os.Getenv("SESSION_KEY")
	if SESSION_KEY == "" {
		log.Fatal("SESSION_KEY not found in env")
	}

	CSRF_KEY = os.Getenv("CSRF_KEY")
	if CSRF_KEY == "" {
		log.Fatal("CSRF_KEY not found in env")
	}
}
