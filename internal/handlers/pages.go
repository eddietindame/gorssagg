package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers/errors"
	"github.com/eddietindame/gorssagg/internal/models"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/eddietindame/gorssagg/internal/templates"
	"github.com/eddietindame/gorssagg/internal/templates/components"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

func handlerWithLayout(contents templ.Component, title string) *templ.ComponentHandler {
	return templ.Handler(templates.Layout(contents, title))
}

func handlerWithDashboardLayout(contents templ.Component, r *http.Request) *templ.ComponentHandler {
	session, _ := store.Store.Get(r, "session")
	username := session.Values["username"].(string)
	return templ.Handler(templates.LayoutFull(templates.DashboardLayout(templates.DashboardProps{
		Contents:    contents,
		CurrentPage: r.URL.Path,
		Username:    username,
	}), "Dashboard"))
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Login(templates.LoginProps{
		CsrfToken:  csrf.Token(r),
		Registered: r.URL.Query().Has("registered"),
		Reset:      r.URL.Query().Has("reset"),
	}), "Login").ServeHTTP(w, r)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Register(csrf.Token(r)), "Register").ServeHTTP(w, r)
}

func ForgotPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Forgot(csrf.Token(r)), "Forgot password").ServeHTTP(w, r)
}

func ResetPageHandler(w http.ResponseWriter, r *http.Request) {
	heading := "Reset password"
	token := r.URL.Query().Get("token")

	if token == "" {
		handlerWithLayout(templates.Reset(templates.ResetProps{
			Err: errors.ResetToken,
		}), heading).ServeHTTP(w, r)
		return
	}

	_, err := store.GetEmailFromToken(r.Context(), token)
	if err != nil {
		handlerWithLayout(templates.Reset(templates.ResetProps{
			Err: errors.ResetToken,
		}), heading).ServeHTTP(w, r)
		return
	}

	handlerWithLayout(templates.Reset(templates.ResetProps{
		CsrfToken:  csrf.Token(r),
		ResetToken: token,
	}), heading).ServeHTTP(w, r)
}

func (apiCfg *APIConfig) PostsPageHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Store.Get(r, "session")
	if err != nil {
		// TODO: redirect / show error
		handlerWithDashboardLayout(components.Posts(components.PostsProps{
			Posts: []models.Post{},
		}), r).ServeHTTP(w, r)
		return
	}

	userId, err := uuid.Parse(session.Values["user_id"].(string))
	if err != nil {
		// TODO: show error
		handlerWithDashboardLayout(components.Posts(components.PostsProps{
			Posts: []models.Post{},
		}), r).ServeHTTP(w, r)
		return
	}

	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: userId,
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
	handlerWithDashboardLayout(components.Feeds(components.FeedsProps{}), r).ServeHTTP(w, r)
}
