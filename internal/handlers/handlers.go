package handlers

import (
	"github.com/eddietindame/gorssagg/internal/database"
)

// APIConfig holds dependencies for handlers
type APIConfig struct {
	DB *database.Queries
}
