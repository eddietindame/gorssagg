package handlers

import (
	"fmt"
	"net/http"

	"github.com/eddietindame/gorssagg/internal/auth"
	"github.com/eddietindame/gorssagg/internal/database"
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
		session, _ := Store.Get(r, "session")
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
