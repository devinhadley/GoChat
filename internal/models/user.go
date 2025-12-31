// Package models handels serialization and database mapping.
package models

import "time"

// User represents a single Row in users table.
// NOTE: We use struct tags to add additional info to the struct that can be used
// with reflection!
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	SignUpDate   *time.Time
	IsActive     bool
}
