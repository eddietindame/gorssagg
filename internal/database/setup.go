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
		db, _ = sql.Open("postgres", config.DB_URL)
		// Open doesn't seem to return an error when there is no db connection
		err = db.Ping()
		if err != nil {
			log.Fatal("Failed to open database connection:", err)
		}
	})
	return db
}

func GetQueries() *Queries {
	return New(getDB())
}
