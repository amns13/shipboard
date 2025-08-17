package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/amns13/shipboard/internal/conf"
	"github.com/amns13/shipboard/internal/services"
)

const AuthUserID = "authenticated_user_id"
const AuthSessionID = "authenticated_session_id"

func RequireAuth(env *conf.Env) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session cookie
			cookie, err := r.Cookie("session_id")
			if err != nil {
				env.Logger.Printf("No session cookie found: %v", err)
				http.Redirect(w, r, "/login/", http.StatusTemporaryRedirect)
				return
			}

			// Validate session using env
			sessionStore := services.SessionStore{
				Client: env.Rdb,
			}

			sessionID := cookie.Value
			sessionData, err := sessionStore.Get(sessionID)
			if err != nil {
				env.Logger.Printf("Invalid session: %v", err)
				http.Redirect(w, r, "/login/", http.StatusTemporaryRedirect)
				return
			}
			if sessionData.ExpiresAt.Before(time.Now()) {
				env.Logger.Printf("Session expired. Logging out")

				err := sessionStore.Expire(sessionID)
				if err != nil {
					env.Logger.Printf("Error expiring session: %v", err)
					http.Redirect(w, r, "/login/", http.StatusTemporaryRedirect)
					return
				}
			}

			// Add the user id to the request context to pass on to further middlewares in the chain
			ctx := context.WithValue(r.Context(), AuthUserID, sessionData.UserID)
			ctx = context.WithValue(ctx, AuthSessionID, sessionID)
			req := r.WithContext(ctx)

			// Session is valid, continue to next handler
			next.ServeHTTP(w, req)
		})
	}
}

