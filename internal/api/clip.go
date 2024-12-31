package api

import (
	"context"
	"log"
	"net/http"

	"github.com/amns13/shipboard/internal/env"
)

func Broadcast(env *env.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		value := req.PostFormValue("value")
		log.Printf("Broadcasting value %s", value)

		ctx := context.Background()
		err := env.Rdb.Set(ctx, "value_1", value, 0).Err()
		if err != nil {
			log.Fatalf("Error while broadcasting clipboard: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

