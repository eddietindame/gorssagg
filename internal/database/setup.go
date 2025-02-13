package database

import (
	"database/sql"
	"log"
	"os"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

func getDB() *sql.DB {
	once.Do(func() {
		dbUrl := os.Getenv("DB_URL")
		if dbUrl == "" {
			log.Fatal("DB_URL not found in env")
		}

		var err error
		db, err = sql.Open("postgres", dbUrl)
		if err != nil {
			log.Fatal("Failed to open database connection:", err)
		}
	})
	return db
}

func GetQueries() *Queries {
	return New(getDB())
}
