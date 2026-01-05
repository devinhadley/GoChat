package store

import (
	"context"
	"time"

	"gochat/main/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionService struct {
	db *pgxpool.Pool
}

func NewSessionService(db *pgxpool.Pool) SessionService {
	return SessionService{
		db: db,
	}
}

func (service *SessionService) CreateSession(ctx context.Context, sessionID string, userID int64, expiresAt time.Time) (models.Session, error) {
	createSessionQuery := `
    INSERT INTO sessions (
        session_id, 
        user_id, 
        expires_at
    ) 
    VALUES ($1, $2, $3)
	  RETURNING *`

	var session models.Session
	err := service.db.QueryRow(ctx, createSessionQuery, sessionID, userID, expiresAt).Scan(
		&session.SessionID,
		&session.UserID,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err != nil {
		return models.Session{}, err
	}
	return session, nil
}
