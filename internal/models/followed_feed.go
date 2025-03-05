package models

import (
	"time"

	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/google/uuid"
)

type FollowedFeed struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Language    string    `json:"language"`
	UserID      uuid.UUID `json:"user_id"`
	FollowID    uuid.UUID `json:"follow_id"`
}

func DatabaseFollowedFeedToFollowedFeed(dbFollowedFeed database.GetFollowedFeedsRow) FollowedFeed {
	return FollowedFeed{
		ID:          dbFollowedFeed.ID,
		CreatedAt:   dbFollowedFeed.CreatedAt,
		UpdatedAt:   dbFollowedFeed.UpdatedAt,
		Name:        dbFollowedFeed.Name,
		Url:         dbFollowedFeed.Url,
		Description: dbFollowedFeed.Description.String,
		Image:       dbFollowedFeed.Image.String,
		Language:    dbFollowedFeed.Language.String,
		UserID:      dbFollowedFeed.UserID,
		FollowID:    dbFollowedFeed.ID_2,
	}
}

func DatabaseFollowedFeedsToFollowedFeeds(dbFollowedFeeds []database.GetFollowedFeedsRow) []FollowedFeed {
	followedFeeds := []FollowedFeed{}
	for _, dbFollowedFeeds := range dbFollowedFeeds {
		followedFeeds = append(followedFeeds, DatabaseFollowedFeedToFollowedFeed(dbFollowedFeeds))
	}
	return followedFeeds
}
