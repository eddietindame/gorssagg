package database

import (
	"database/sql"
	"log"
	"sync"

	"github.com/eddietindame/gorssagg/internal/config"
)

var (
	db   *sql.DB
	once sync.Once
)

func getDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", config.DB_URL)
		if err != nil {
			log.Fatal("Failed to open database connection:", err)
		}
	})
	return db
}

func GetQueries() *Queries {
	return New(getDB())
}
