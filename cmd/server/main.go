package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amns13/shipboard/internal/api"
	"github.com/amns13/shipboard/internal/env"
	"github.com/joho/godotenv"
)

func startServer() {
	s := &http.Server{
		Addr: ":8080",
	}
	log.Fatal(s.ListenAndServe())
}

func registerEndpoints(env *env.Env) {
	http.Handle("/", http.NotFoundHandler())

	http.HandleFunc("POST /clip/", api.Broadcast(env))
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	redisUri := os.Getenv("REDIS_URI")

	env, err := env.GetEnv(redisUri)
	if err != nil {
		log.Fatalf("Error initializing environment: %v", err)
	}

	registerEndpoints(env)
	startServer()
}
