package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers/ctx"
	"github.com/eddietindame/gorssagg/internal/handlers/errors"
	"github.com/eddietindame/gorssagg/internal/models"
	"github.com/eddietindame/gorssagg/internal/rss"
	"github.com/eddietindame/gorssagg/internal/templates/components"
	"github.com/google/uuid"
)

func (apiCfg *APIConfig) CreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	ts := time.Now().UTC()

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: ts,
		UpdatedAt: ts,
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}

	respondWithJSON(w, 201, models.DatabaseFeedToFeed(feed))
}

// func handlerWithFollowedFeeds(csrfToken string, feedUrl string, followedFeeds []database.FeedFollow) *templ.ComponentHandler {
// 	templ.Handler(
// 		components.FollowedFeeds(components.FollowedFeedsProps{
// 			CsrfToken:     csrf.Token(r),
// 			FeedUrl:       url,
// 			FollowedFeeds: models.DatabaseFollowedFeedsToFollowedFeeds(feeds),
// 		}),
// 	)
// }

func (apiCfg *APIConfig) FeedHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("feed_url")

	user, ok := ctx.GetUserFromContext(r.Context())
	if !ok {
		// TODO: handle error / redirect
		log.Println("User not in context")
		return
	}

	ok, _ = regexp.MatchString("^https?:\\/\\/[^\\s\\/$.?#].[^\\s]*\\.xml(\\?[^\\s]*)?(#[^\\s]*)?$", url)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		templ.Handler(
			components.FormError(components.FormErrorProps{
				Error: errors.FeedInvalid,
			}),
		).ServeHTTP(w, r)
		return
	}

	// Attempt to get an existing database feed
	feed, err := apiCfg.DB.GetFeedByUrl(r.Context(), url)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			// If none exists, create one
			rssFeed, err := rss.UrlToFeed(url) // TODO: handle malformed XML / HTML
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				templ.Handler(
					components.FormError(components.FormErrorProps{
						Error: errors.FeedFetch,
					}),
				).ServeHTTP(w, r)
				return
			}

			ts := time.Now().UTC()

			description := sql.NullString{}
			if rssFeed.Channel.Description != "" {
				description.String = rssFeed.Channel.Description
				description.Valid = true
			}

			image := sql.NullString{}
			if rssFeed.Channel.Image.Url != "" {
				image.String = rssFeed.Channel.Image.Url
				image.Valid = true
			}

			language := sql.NullString{}
			if rssFeed.Channel.Language != "" {
				language.String = rssFeed.Channel.Language
				language.Valid = true
			}

			feed, err = apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
				ID:          uuid.New(),
				CreatedAt:   ts,
				UpdatedAt:   ts,
				Name:        rssFeed.Channel.Title,
				Description: description,
				Language:    language,
				Image:       image,
				Url:         url,
				UserID:      user.UserID,
			})
			if err != nil {
				log.Println("error creating feed", err)
				w.WriteHeader(http.StatusInternalServerError)
				templ.Handler(
					components.FormError(components.FormErrorProps{
						Error: errors.FeedCreate,
					}),
				).ServeHTTP(w, r)
				return
			}
		} else {
			log.Println("error fetching database feed", err)
			w.WriteHeader(http.StatusInternalServerError)
			templ.Handler(
				components.FormError(components.FormErrorProps{
					Error: errors.FeedFetch,
				}),
			).ServeHTTP(w, r)
			return
		}
	}

	ts := time.Now().UTC()

	// Then follow the feed
	_, err = apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: ts,
		UpdatedAt: ts,
		UserID:    user.UserID,
		FeedID:    feed.ID,
	})
	if err != nil {
		log.Println("error creating feed_follow", err)
		w.WriteHeader(http.StatusInternalServerError)
		templ.Handler(
			components.FormError(components.FormErrorProps{
				Error: errors.FollowCreate,
			}),
		).ServeHTTP(w, r)
		return
	}

	// And list followed feeds
	feeds, err := apiCfg.DB.GetFollowedFeeds(r.Context(), user.UserID)
	if err != nil {
		log.Println("error fetching followed feeds", err)
		w.WriteHeader(http.StatusInternalServerError)
		templ.Handler(
			components.FormError(components.FormErrorProps{
				Error: errors.FollowedFeedsRead,
			}),
		).ServeHTTP(w, r)
		return
	}

	templ.Handler(
		components.FollowedFeedList(components.FollowedFeedListProps{
			FollowedFeeds: models.DatabaseFollowedFeedsToFollowedFeeds(feeds),
		}),
	).ServeHTTP(w, r)
}

func (apiCfg *APIConfig) GetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feeds: %v", err))
		return
	}

	respondWithJSON(w, 200, models.DatabaseFeedsToFeeds(feeds))
}
