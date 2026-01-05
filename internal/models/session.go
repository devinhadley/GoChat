package models

import "time"

// TODO: Move me to store package (in sessions)...

type Session struct {
	SessionID string
	UserID    int64
	ExpiresAt time.Time
	CreatedAt time.Time
}
