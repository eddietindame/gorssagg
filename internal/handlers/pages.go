package handlers

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers/ctx"
	"github.com/eddietindame/gorssagg/internal/handlers/errors"
	"github.com/eddietindame/gorssagg/internal/models"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/eddietindame/gorssagg/internal/templates"
	"github.com/eddietindame/gorssagg/internal/templates/components"
	"github.com/gorilla/csrf"
)

func handlerWithLayout(contents templ.Component, title string, csrfToken string) *templ.ComponentHandler {
	return templ.Handler(templates.Layout(contents, title, csrfToken))
}

func handlerWithDashboardLayout(contents templ.Component, r *http.Request) *templ.ComponentHandler {
	session, _ := store.Store.Get(r, "session")
	username := session.Values["username"].(string)
	return templ.Handler(templates.LayoutFull(templates.DashboardLayout(templates.DashboardProps{
		Contents:    contents,
		CurrentPage: r.URL.Path,
		Username:    username,
	}), "Dashboard", csrf.Token(r)))
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	handlerWithLayout(templates.Login(templates.LoginProps{
		CsrfToken:  csrfToken,
		Registered: r.URL.Query().Has("registered"),
		Reset:      r.URL.Query().Has("reset"),
	}), "Login", csrfToken).ServeHTTP(w, r)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	handlerWithLayout(templates.Register(csrfToken), "Register", csrfToken).ServeHTTP(w, r)
}

func ForgotPageHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	handlerWithLayout(templates.Forgot(csrfToken), "Forgot password", csrfToken).ServeHTTP(w, r)
}

func ResetPageHandler(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)
	heading := "Reset password"
	token := r.URL.Query().Get("token")

	if token == "" {
		handlerWithLayout(templates.Reset(templates.ResetProps{
			Err: errors.ResetToken,
		}), heading, csrfToken).ServeHTTP(w, r)
		return
	}

	_, err := store.GetEmailFromToken(r.Context(), token)
	if err != nil {
		handlerWithLayout(templates.Reset(templates.ResetProps{
			Err: errors.ResetToken,
		}), heading, csrfToken).ServeHTTP(w, r)
		return
	}

	handlerWithLayout(templates.Reset(templates.ResetProps{
		CsrfToken:  csrf.Token(r),
		ResetToken: token,
	}), heading, csrfToken).ServeHTTP(w, r)
}

func (apiCfg *APIConfig) PostsPageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := ctx.GetUserFromContext(r.Context())
	if !ok {
		// TODO: handle error / redirect
		log.Println("User not in context")
		return
	}

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.UserID,
	})
	if err != nil {
		// TODO: show error
		handlerWithDashboardLayout(components.Posts(components.PostsProps{
			Posts: []models.Post{},
		}), r).ServeHTTP(w, r)
		return
	}

	handlerWithDashboardLayout(components.Posts(components.PostsProps{
		Posts: models.DatabasePostsToPosts(posts),
	}), r).ServeHTTP(w, r)
}

func (apiCfg *APIConfig) FeedsPageHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := ctx.GetUserFromContext(r.Context())
	if !ok {
		// TODO: handle error / redirect
		log.Println("User not in context")
		return
	}

	// TODO: handle session errors
	feeds, err := apiCfg.DB.GetFollowedFeeds(r.Context(), user.UserID)
	if err != nil {
		// TODO: show error
		handlerWithDashboardLayout(components.FollowedFeeds(components.FollowedFeedsProps{
			CsrfToken: csrf.Token(r),
		}), r).ServeHTTP(w, r)
		return
	}

	handlerWithDashboardLayout(components.FollowedFeeds(components.FollowedFeedsProps{
		FollowedFeeds: models.DatabaseFollowedFeedsToFollowedFeeds(feeds),
		CsrfToken:     csrf.Token(r),
	}), r).ServeHTTP(w, r)
}
