// Package passwords provides functions related to passwords & password hashing.
package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2IdParams struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

// DefaultArgon2Params should reflect OSWAP minimum recs: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
var DefaultArgon2Params = Argon2IdParams{
	Time:    2,
	Memory:  19 * 1024,
	Threads: 1,
}

func DoesPasswordMatchHashedPassword(password string, hashString string) (bool, error) {
	parts := strings.Split(hashString, "$")

	if len(parts) != 6 {
		return false, errors.New("hash string does not have six parts")
	}

	// Skip the version.

	// Time, Memory, & Threads
	params := Argon2IdParams{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Time, &params.Threads)
	if err != nil {
		return false, errors.New("failed to extract time memory and threads from hash")
	}

	// Salt
	existingSalt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, errors.New("failed to base64 decode the salt")
	}

	// Hash
	existingHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, errors.New("failed to base64 decode the hash")
	}

	hashedPass := argon2.IDKey(
		[]byte(password),
		existingSalt,
		params.Time,
		params.Memory,
		params.Threads,
		uint32(len(existingHash)))

	return subtle.ConstantTimeCompare(hashedPass, existingHash) == 1, nil
}

func CreatePasswordHash(password string, params Argon2IdParams) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hashBytes := argon2.IDKey(
		[]byte(password),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		32)

	hashAsBase64String := base64.RawStdEncoding.EncodeToString(hashBytes)
	saltAsBase64String := base64.RawStdEncoding.EncodeToString(salt)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, params.Memory, params.Time, params.Threads, saltAsBase64String, hashAsBase64String)

	return encodedHash, nil
}

// generateSalt creates a 16-byte cryptographically secure random salt.
func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}
