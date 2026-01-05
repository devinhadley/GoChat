// Package sessions provides some utilities for creating user sessions.
package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

const SessionCookieName string = "goChatSessionID"

type ContextKey string

const (
	UserContextKey ContextKey = "User"
)

func CreateSessionCookie() (http.Cookie, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return http.Cookie{}, nil
	}
	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)

	return http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Secure:   true,
		HttpOnly: true,
		Expires:  thirtyDaysFromNow,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}, nil
}

func CreateClearSessionCookie() http.Cookie {
	return http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
}

// GenerateSessionID generates a random byte array with 128 bits of entropy
// and returns it as a base64 encoded string.
func generateSessionID() (string, error) {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(bytes), nil
}
