package api

import (
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/amns13/shipboard/internal/conf"
	"github.com/amns13/shipboard/internal/middleware"
	"github.com/amns13/shipboard/internal/model"
	"github.com/amns13/shipboard/internal/services"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const INTERNAL_SERVER_ERROR = "An error occurred. Please contact support."

func Register(env *conf.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email, err := mail.ParseAddress(req.PostFormValue("email"))
		if err != nil {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		exists, err := model.UserExists(env, email)
		if err != nil {
			log.Printf("Error occurred: %v", err)
			http.Error(w, INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}
		if exists {
			log.Printf("Email already exists: %v", email)
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}

		// TODO: Input validations for name and password length
		// TODO: Password strength validation
		name := req.PostFormValue("name")
		password := req.PostFormValue("password")
		// TODO: Randomly giving cost 14. Confirm an optimal value
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
		if err != nil {
			log.Printf("Error occurred while generating password hash: %v", err)
			http.Error(w, INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}

		userData := model.UserCreator{
			Name:         name,
			Email:        email.Address,
			PasswordHash: string(passwordHash),
		}
		_, err = userData.Create(env)
		if err != nil {
			log.Printf("Error occurred while creating user: %v", err)
			http.Error(w, INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func RegistrationForm(env *conf.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := env.Templates.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			log.Printf("Error occurred while rendering registration form: %v", err)
			http.Error(w, INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}
	}
}

func Login(env *conf.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email, err := mail.ParseAddress(req.PostFormValue("email"))
		if err != nil {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		user, err := model.GetUserByEmail(env, email.Address)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Printf("Email not found: %v", email.Address)
				http.Error(w, "Invalid email or password", http.StatusBadRequest)
			} else {
				log.Printf("Error occurred while fetching user: %v", err)
				http.Error(w, INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			}
			return
		}

		password := req.PostFormValue("password")
		// TODO: Randomly giving cost 14. Confirm an optimal value
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusBadRequest)
			return
		}

		sessionData := services.SessionData{
			UserID:    user.Id,
			Email:     user.Email,
			LoginTime: time.Now(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		sessionStore := services.SessionStore{
			Client: env.Rdb,
		}
		sessionID, err := sessionStore.CreateSession(sessionData)
		if err != nil {
			log.Printf("Error occurred while creating user session: %v", err)
			http.Error(w, INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  sessionData.ExpiresAt,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.Header().Set("HX-Redirect", "/clip/")
		w.WriteHeader(http.StatusOK)
	}
}

func LoginForm(env *conf.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := env.Templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Printf("Error occurred while rendering login form: %v", err)
			http.Error(w, INTERNAL_SERVER_ERROR, http.StatusInternalServerError)
			return
		}
	}
}

func Logout(env *conf.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessionID, ok := req.Context().Value(middleware.AuthSessionID).(string)
		if !ok {
			log.Println("Invalid session id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionStore := services.SessionStore{
			Client: env.Rdb,
		}
		err := sessionStore.Expire(sessionID)
		if err != nil {
			log.Printf("Error occurred while expiring session %s: %v", sessionID, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Set cookie with past expiration to delete it
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Expires:  time.Unix(0, 0), // Past date
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})
		w.WriteHeader(http.StatusNoContent)
	}
}
