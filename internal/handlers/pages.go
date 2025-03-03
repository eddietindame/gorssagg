package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/eddietindame/gorssagg/internal/templates"
	"github.com/gorilla/csrf"
)

func handlerWithLayout(contents templ.Component, title string) *templ.ComponentHandler {
	return templ.Handler(templates.Layout(contents, title))
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Login(csrf.Token(r)), "Login").ServeHTTP(w, r)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Register(csrf.Token(r)), "Register").ServeHTTP(w, r)
}

func ForgotPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Forgot(csrf.Token(r)), "Forgot password").ServeHTTP(w, r)
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

	handlerWithLayout(templates.Reset(csrf.Token(r), token), "Reset password").ServeHTTP(w, r)
}

func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	handlerWithLayout(templates.Dashboard(), "Dashboard").ServeHTTP(w, r)
}
