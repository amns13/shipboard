package api

import (
	"context"
	"net/http"
	// "time"

	"github.com/amns13/shipboard/internal/conf"
)

func Clip(env *conf.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := env.Templates.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Broadcast(env *conf.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// time.Sleep(1 * time.Second)
		value := req.PostFormValue("content")
		// Testing
		env.Logger.Printf("Broadcasting value %s", value)

		ctx := context.Background()
		err := env.Rdb.Set(ctx, "value_1", value, 0).Err()
		if err != nil {
			env.Logger.Printf("Error while broadcasting clipboard: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
