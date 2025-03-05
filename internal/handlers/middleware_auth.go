package handlers

import (
	"fmt"
	"net/http"

	"github.com/eddietindame/gorssagg/internal/auth"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers/ctx"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/google/uuid"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *APIConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}

		handler(w, r, user)
	}
}

func MiddlewareSessionAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Store.Get(r, "session")
		if err != nil {
			http.Error(w, "Session error", http.StatusInternalServerError)
			return
		}

		auth, ok := session.Values[store.Authenticated].(bool)
		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		email := session.Values[store.Email].(string)
		username := session.Values[store.Username].(string)
		userId, err := uuid.Parse(session.Values[store.UserID].(string))
		if err != nil || email == "" || username == "" {
			// TODO: redirect
			return
		}

		user := ctx.UserContext{
			UserID:   userId,
			Username: username,
			Email:    email,
		}

		handler.ServeHTTP(w, r.WithContext(ctx.NewContextWithUser(r.Context(), user)))
	})
}
