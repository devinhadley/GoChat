// Package middleware contains the functions which run before the mutex handles request routing.
// In other words, they run before every request is handeled.
package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"

	"gochat/main/internal/store"
	"gochat/main/internal/utils/sessions"

	"github.com/jackc/pgx/v5"
)

// AuthMiddleware populates the User struct if the request contains a valid session id.
func AuthMiddleware(next http.Handler, userService store.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(sessions.SessionCookieName)
		if err != nil {
			if !errors.Is(err, http.ErrNoCookie) {
				log.Printf("Error when retrieving session id from cookie: %v", err)
			}

			next.ServeHTTP(w, r)
			return
		}

		user, err := userService.GetUserFromSessionID(r.Context(), sessionCookie.Value)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// The session is likely invalidated or expired.
				clearSessionCookie := sessions.CreateClearSessionCookie()
				http.SetCookie(w, &clearSessionCookie)
			} else {
				log.Printf("Error when getting user from session: %v", err)
			}
			next.ServeHTTP(w, r)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), sessions.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}
