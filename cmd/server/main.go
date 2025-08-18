// Package main
// This file contains the code for starting the http server.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amns13/shipboard/internal/api"
	"github.com/amns13/shipboard/internal/conf"
	"github.com/amns13/shipboard/internal/middleware"
	"github.com/joho/godotenv"
)

var templates = []string{"templates/index.html", "templates/register.html", "templates/login.html"}

func startServer(mux *http.ServeMux) {
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(s.ListenAndServe())
}

func registerEndpoints(mux *http.ServeMux, env *conf.Env) {
	// Create middleware functions
	requestMiddleware := middleware.LogRequestResponse(env)
	authMiddleware := middleware.RequireAuth(env)

	// Restrict root path
	mux.Handle("/", http.NotFoundHandler())

	// Static files (no middleware needed)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Public routes with logging only
	mux.Handle("GET /register/", requestMiddleware(http.HandlerFunc(api.RegistrationForm(env))))
	mux.Handle("POST /register/", requestMiddleware(http.HandlerFunc(api.Register(env))))
	mux.Handle("GET /login/", requestMiddleware(http.HandlerFunc(api.LoginForm(env))))
	mux.Handle("POST /login/", requestMiddleware(http.HandlerFunc(api.Login(env))))

	// Protected routes with both logging and auth

	mux.Handle("DELETE /logout/", requestMiddleware(authMiddleware(http.HandlerFunc(api.Logout(env)))))
	mux.Handle("GET /clip/", requestMiddleware(authMiddleware(http.HandlerFunc(api.Clip(env)))))
	mux.Handle("POST /clip/", requestMiddleware(authMiddleware(http.HandlerFunc(api.Broadcast(env)))))
}

func loadEnvironment() (*conf.Env, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	env, err := conf.LoadEnv(os.Getenv("POSTGRES_URI"), os.Getenv("REDIS_URI"), templates)
	if err != nil {
		return nil, err
	}
	return env, err
}

func main() {
	env, err := loadEnvironment()
	if err != nil {
		log.Fatalf("Error initializing environment: %v", err)
	}
	env.Logger.Println("Initialized environment")
	defer env.Db.Close()

	mux := http.NewServeMux()
	registerEndpoints(mux, env)
	startServer(mux)
}
