package api

import (
	"context"
	"fmt"
	"net/http"

	// "time"

	"github.com/amns13/shipboard/internal/conf"
	"github.com/amns13/shipboard/internal/middleware"
	"github.com/amns13/shipboard/internal/model"
)

const clipboardPrefix = "__clip__"

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
		userID, ok := req.Context().Value(middleware.AuthUserID).(int32)
		if !ok {
			env.Logger.Println("Invalid user id", userID)
			http.Redirect(w, req, "/logout/", http.StatusTemporaryRedirect)
			return
		}
		value := req.PostFormValue("content")
		user, err := model.GetUserByID(env, userID)
		if err != nil {
			env.Logger.Println("Invalid user id", userID)
			http.Redirect(w, req, "/logout/", http.StatusTemporaryRedirect)
			return
		}
		key := fmt.Sprintf("%s%s", clipboardPrefix, user.Uid)

		ctx := context.Background()
		//TODO: Add a default TTL
		err = env.Rdb.Set(ctx, key, value, 0).Err()
		if err != nil {
			env.Logger.Printf("Error while broadcasting clipboard: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// TODO: This should return 201 created. Bu, fsr the input box is removed
		// on that. Debug and fix
		w.WriteHeader(http.StatusNoContent)
	}
}
