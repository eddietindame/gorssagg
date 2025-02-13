package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/templates"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.Login()).ServeHTTP(w, r)
}

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.Register()).ServeHTTP(w, r)
}

func DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	templ.Handler(templates.Dashboard()).ServeHTTP(w, r)
}
