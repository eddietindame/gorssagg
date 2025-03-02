package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/eddietindame/gorssagg/internal/templates"
	"github.com/gorilla/csrf"
)

var csrfFormKey = "gorilla.csrf.Token"

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(
		templates.Layout(templates.Login(csrf.Token(r), csrfFormKey), "Login"),
	).ServeHTTP(w, r)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(
		templates.Layout(templates.Register(csrf.Token(r), csrfFormKey), "Register"),
	).ServeHTTP(w, r)
}

func ForgotPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(
		templates.Layout(templates.Forgot(csrf.Token(r), csrfFormKey), "Forgot password"),
	).ServeHTTP(w, r)
}

func ResetPageHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	_, err := store.GetEmailFromToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	templ.Handler(
		templates.Layout(templates.Reset(csrf.Token(r), csrfFormKey, token), "Reset password"),
	).ServeHTTP(w, r)
}

func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(
		templates.Layout(templates.Dashboard(), "Dashboard"),
	).ServeHTTP(w, r)
}
