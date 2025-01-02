// Package main
// This file contains the code for starting the http server.
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

func loadEnvironent() (*env.Env, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	loadedEnv, err := env.LoadEnv(os.Getenv("POSTGRES_URI"), os.Getenv("REDIS_URI"))
	if err != nil {
		return nil, err
	}
	return loadedEnv, err
}

func main() {
	loadedEnv, err := loadEnvironent()
	if err != nil {
		log.Fatalf("Error initializing environment: %v", err)
	}
	defer loadedEnv.Db.Close()

	registerEndpoints(loadedEnv)
	startServer()
}
