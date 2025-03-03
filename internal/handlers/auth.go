package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/eddietindame/gorssagg/internal/config"
	"github.com/eddietindame/gorssagg/internal/database"
	"github.com/eddietindame/gorssagg/internal/handlers/errors"
	"github.com/eddietindame/gorssagg/internal/mail"
	"github.com/eddietindame/gorssagg/internal/store"
	"github.com/eddietindame/gorssagg/internal/templates/components"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

func responseWithLoginForm(csrfToken string, values components.LoginFormValues, err errors.HandlerError) *templ.ComponentHandler {
	return templ.Handler(components.LoginForm(components.LoginFormProps{
		CsrfToken: csrfToken,
		Err:       err,
		Values:    values,
	}))
}

func (apiCfg *APIConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	rememberMeStr := r.FormValue("remember_me")
	rememberMe, err := strconv.ParseBool(rememberMeStr)
	if err != nil {
		rememberMe = false
	}

	values := components.LoginFormValues{
		Username:   username,
		Password:   password,
		RememberMe: rememberMe,
	}

	hashedPassword, err := apiCfg.DB.GetUserPassword(r.Context(), username)
	if err != nil {
		if r.Header.Get("Hx-Request") != "" {
			responseWithLoginForm(csrf.Token(r), values, errors.LoginCredentials).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.LoginCredentials.ToString(), http.StatusUnauthorized)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if r.Header.Get("Hx-Request") != "" {
			responseWithLoginForm(csrf.Token(r), values, errors.LoginCredentials).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.LoginCredentials.ToString(), http.StatusUnauthorized)
		}
		return
	}

	session, err := store.Store.Get(r, "session")
	if err != nil {
		log.Println(err)
		if r.Header.Get("Hx-Request") != "" {
			responseWithLoginForm(csrf.Token(r), values, errors.SessionError).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.SessionError.ToString(), http.StatusInternalServerError)
		}
		return
	}

	session.Values["authenticated"] = true
	session.Values["username"] = username

	if rememberMe {
		session.Options.MaxAge = 86400 * 30 // 30 days
	} else {
		session.Options.MaxAge = 3600 // 1 hour
	}

	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
		if r.Header.Get("Hx-Request") != "" {
			responseWithLoginForm(csrf.Token(r), values, errors.SessionError).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.SessionError.ToString(), http.StatusInternalServerError)
		}
		return
	}

	redirect(w, r, "/dashboard")
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

func responseWithRegisterForm(csrfToken string, values components.RegisterFormValues, err errors.HandlerError) *templ.ComponentHandler {
	return templ.Handler(components.RegisterForm(components.RegisterFormProps{
		CsrfToken: csrfToken,
		Err:       err,
		Values:    values,
	}))
}

func (apiCfg *APIConfig) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirm := r.FormValue("password_confirm")

	values := components.RegisterFormValues{
		Username:        username,
		Email:           email,
		Password:        password,
		ConfirmPassword: passwordConfirm,
	}

	if password != passwordConfirm {
		if r.Header.Get("Hx-Request") != "" {
			responseWithRegisterForm(csrf.Token(r), values, errors.RegisterPassword).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.RegisterPassword.ToString(), http.StatusBadRequest)
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		if r.Header.Get("Hx-Request") != "" {
			responseWithRegisterForm(csrf.Token(r), values, errors.ServerError).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.ServerError.ToString(), http.StatusInternalServerError)
		}
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
			if r.Header.Get("Hx-Request") != "" {
				responseWithRegisterForm(csrf.Token(r), values, errors.RegisterUsername).ServeHTTP(w, r)
			} else {
				http.Error(w, errors.RegisterUsername.ToString(), http.StatusBadRequest)
			}
		} else if strings.Contains(err.Error(), "email_address_check") {
			if r.Header.Get("Hx-Request") != "" {
				responseWithRegisterForm(csrf.Token(r), values, errors.RegisterEmail).ServeHTTP(w, r)
			} else {
				http.Error(w, errors.RegisterEmail.ToString(), http.StatusBadRequest)
			}
		} else {
			if r.Header.Get("Hx-Request") != "" {
				responseWithRegisterForm(csrf.Token(r), values, errors.RegisterUserExists).ServeHTTP(w, r)
			} else {
				http.Error(w, errors.RegisterUserExists.ToString(), http.StatusBadRequest)
			}
		}
		return
	}

	redirect(w, r, "/login")
}

func responseWithForgotForm(csrfToken string, values components.ForgotFormValues, err errors.HandlerError) *templ.ComponentHandler {
	return templ.Handler(components.ForgotForm(components.ForgotFormProps{
		CsrfToken: csrfToken,
		Err:       err,
		Values:    values,
	}))
}

func (apiCfg *APIConfig) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	values := components.ForgotFormValues{
		Email: email,
	}

	_, err := apiCfg.DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		log.Println("Error finding user:", err)
		if r.Header.Get("Hx-Request") != "" {
			responseWithForgotForm(csrf.Token(r), values, errors.ForgotNotFound).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.ForgotNotFound.ToString(), http.StatusNotFound)
		}
		return
	}

	token := uuid.New().String()

	err = store.StoreToken(r.Context(), token, email, 30*time.Minute)
	if err != nil {
		log.Println("Error generating reset token:", err)
		if r.Header.Get("Hx-Request") != "" {
			responseWithForgotForm(csrf.Token(r), values, errors.ForgotToken).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.ForgotToken.ToString(), http.StatusInternalServerError)
		}
		return
	}

	resetLink := fmt.Sprintf("http://%s/reset-password?token=%s", config.HOST, token)
	err = mail.SendResetEmail(email, resetLink)
	if err != nil {
		store.DeleteToken(r.Context(), token)
		log.Println("Error sending reset email:", err)
		if r.Header.Get("Hx-Request") != "" {
			responseWithForgotForm(csrf.Token(r), values, errors.ForgotSend).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.ForgotSend.ToString(), http.StatusInternalServerError)
		}
		return
	}

	if r.Header.Get("Hx-Request") != "" {
		templ.Handler(components.ForgotForm(components.ForgotFormProps{
			Success: true,
		})).ServeHTTP(w, r)
	} else {
		w.Write([]byte("Password reset link sent! Check your email."))
	}
}

func responseWithResetForm(csrfToken string, values components.ResetFormValues, err errors.HandlerError) *templ.ComponentHandler {
	return templ.Handler(components.ResetForm(components.ResetFormProps{
		CsrfToken: csrfToken,
		Err:       err,
		Values:    values,
	}))
}

func (apiCfg *APIConfig) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	newPassword := r.FormValue("password")
	newPasswordConfirm := r.FormValue("password_confirm")

	values := components.ResetFormValues{
		Password:        newPassword,
		PasswordConfirm: newPasswordConfirm,
	}

	if newPassword != newPasswordConfirm {
		if r.Header.Get("Hx-Request") != "" {
			responseWithResetForm(csrf.Token(r), values, errors.ResetPassword).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.ResetPassword.ToString(), http.StatusBadRequest)
		}
		return
	}

	email, err := store.GetEmailFromToken(r.Context(), token)
	if err != nil {
		if r.Header.Get("Hx-Request") != "" {
			responseWithResetForm(csrf.Token(r), values, errors.ResetToken).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.ResetToken.ToString(), http.StatusBadRequest)
		}
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	err = apiCfg.DB.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		if r.Header.Get("Hx-Request") != "" {
			responseWithResetForm(csrf.Token(r), values, errors.ResetFailed).ServeHTTP(w, r)
		} else {
			http.Error(w, errors.ResetFailed.ToString(), http.StatusInternalServerError)
		}
		return
	}

	store.DeleteToken(r.Context(), token)

	redirect(w, r, "/login?reset")
}
