package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/handlers/errors"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/eddietindame/gorssagg/internal/templates"
	"github.com/gorilla/csrf"
)

func handlerWithLayout(contents templ.Component, title string) *templ.ComponentHandler {
	return templ.Handler(templates.Layout(contents, title))
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Login(templates.LoginProps{
		CsrfToken: csrf.Token(r),
		Reset:     r.URL.Query().Has("reset"),
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

func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Dashboard(), "Dashboard").ServeHTTP(w, r)
}
