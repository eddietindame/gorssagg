package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/mail"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

	session, err := store.Store.Get(r, "session")
	if err != nil {
		log.Println(err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["username"] = username

	if r.FormValue("remember_me") == "true" {
		session.Options.MaxAge = 86400 * 30 // 30 days
	} else {
		session.Options.MaxAge = 3600 // 1 hour
	}

	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1 // Expire session immediately
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (apiCfg *APIConfig) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	if password != passwordConfirm {
		http.Error(w, "Password does not match confirmation", http.StatusBadRequest)
		return
	}

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
		Email:     email,
		Password:  string(hashedPassword),
	})
	if err != nil {
		log.Println("Error creating user:", err)
		if strings.Contains(err.Error(), "username_check") {
			http.Error(w, "Invalid username", http.StatusBadRequest)
		} else {
			http.Error(w, "User already exists", http.StatusBadRequest)
		}
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (apiCfg *APIConfig) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	_, err := apiCfg.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		log.Println("Error finding user:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	token := uuid.New().String()

	err = store.StoreToken(r.Context(), token, email, 30*time.Minute)
	if err != nil {
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}

	resetLink := fmt.Sprintf("http://%s/reset-password?token=%s", config.HOST, token)
	err = mail.SendResetEmail(email, resetLink)
	if err != nil {
		store.DeleteToken(r.Context(), token)
		log.Println("Error sending reset email:", err)
		http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Password reset link sent! Check your email."))
}

func (apiCfg *APIConfig) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	newPassword := r.FormValue("password")
	newPasswordConfirm := r.FormValue("password_confirm")

	if newPassword != newPasswordConfirm {
		http.Error(w, "Password does not match confirmation", http.StatusBadRequest)
		return
	}

	email, err := store.GetEmailFromToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	err = apiCfg.DB.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}

	store.DeleteToken(r.Context(), token)

	w.Write([]byte("Password reset successful! You can now log in."))
}
