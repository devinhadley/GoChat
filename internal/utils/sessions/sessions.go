// Package sessions provides some utilities for creating user sessions.
package sessions

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateSessionID generates a random byte array with 128 bits of entropy
// and returns it as a base64 encoded string.
func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", nil
	}

	return base64.RawStdEncoding.EncodeToString(bytes), nil
}
