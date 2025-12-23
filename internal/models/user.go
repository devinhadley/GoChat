// Package models handels serialization and database mapping.
package models

import "time"

// User represents a single Row in users table.
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	SignUpDate   *time.Time
	IsActive     bool
}
