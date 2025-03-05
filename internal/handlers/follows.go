package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers/ctx"
	"github.com/eddietindame/gorssagg/internal/models"
	"github.com/eddietindame/gorssagg/internal/templates/components"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var FeedFollowsIdPageParam = "feedFollowID"

func (apiCfg *APIConfig) CreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ts := time.Now().UTC()

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: ts,
		UpdatedAt: ts,
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, models.DatabaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *APIConfig) GetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feed follows: %v", err))
		return
	}
	respondWithJSON(w, 200, models.DatabaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *APIConfig) DeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing feed follow id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follow: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}

func (apiCfg *APIConfig) DeleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	feedFollowIDStr := chi.URLParam(r, FeedFollowsIdPageParam)
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		// TODO: respond with html
		respondWithError(w, 400, fmt.Sprintf("Error parsing feed follow id: %v", err))
		return
	}

	user, ok := ctx.GetUserFromContext(r.Context())
	if !ok {
		// TODO: handle error / redirect
		log.Println("User not in context")
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.UserID,
	})
	if err != nil {
		// TODO: respond with html
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follow: %v", err))
		return
	}

	// And list followed feeds
	feeds, err := apiCfg.DB.GetFollowedFeeds(r.Context(), user.UserID)
	if err != nil {
		log.Println("error creating feed_follow", err)
		// TODO: respond with html
		respondWithError(w, 400, fmt.Sprintf("Error fetching followed feeds: %v", err))
		return
	}

	templ.Handler(
		components.FollowedFeedList(components.FollowedFeedListProps{
			FollowedFeeds: models.DatabaseFollowedFeedsToFollowedFeeds(feeds),
		}),
	).ServeHTTP(w, r)
}
