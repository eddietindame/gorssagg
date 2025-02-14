package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/templates"
	"github.com/gorilla/csrf"
)

var csrfFormKey = "gorilla.csrf.Token"

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.Login(csrf.Token(r), csrfFormKey)).ServeHTTP(w, r)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.Register(csrf.Token(r), csrfFormKey)).ServeHTTP(w, r)
}

func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.Dashboard()).ServeHTTP(w, r)
}
