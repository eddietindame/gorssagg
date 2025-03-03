package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/models"
	"github.com/google/uuid"
)

func (apiCfg *APIConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ts := time.Now().UTC()

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: ts,
		UpdatedAt: ts,
		Username:  params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, 201, models.DatabaseUserToUser(user))
}

func (apiCfg *APIConfig) GetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 201, models.DatabaseUserToUser(user))
}

func (apiCfg *APIConfig) GetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting posts for user: %v", err))
		return
	}

	respondWithJSON(w, 200, models.DatabasePostsToPosts(posts))
}
