package middleware

import (
	"net/http"
	"time"

	"github.com/amns13/shipboard/internal/conf"
)

func LogRequestResponse(env *conf.Env) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			env.Logger.Println(r.Method, r.URL.Path, time.Since(start))
		})
	}
}
