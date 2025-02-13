package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var Store *sessions.CookieStore

func (apiCfg *APIConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPassword, err := apiCfg.DB.GetUserPassword(r.Context(), username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, err := Store.Get(r, "session")
	if err != nil {
		log.Println(err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values["username"] = username

	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	delete(session.Values, "username")
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (apiCfg *APIConfig) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	ts := time.Now().UTC()

	_, err = apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: ts,
		UpdatedAt: ts,
		Username:  username,
		Password:  string(hashedPassword),
	})
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
