// Package store provides the service layer exposing databse operations.
package store

import (
	"context"
	"errors"
	"time"

	"gochat/main/internal/utils/passwords"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) UserService {
	return UserService{
		db: db,
	}
}

type User struct {
	ID           int64
	Username     string
	passwordHash string
	SignUpDate   *time.Time
	IsActive     bool
}

func (store *UserService) CreateUser(username string, password string, context context.Context) (User, error) {
	passHash, err := passwords.CreatePasswordHash(password, passwords.DefaultArgon2Params)
	if err != nil {
		return User{}, err
	}

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING *"

	var user User
	err = store.db.QueryRow(context, query, username, passHash).Scan(
		&user.ID,
		&user.Username,
		&user.passwordHash,
		&user.SignUpDate,
		&user.IsActive,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

var ErrInvalidCredentials = errors.New("no user with the following credentials found")

func (store *UserService) AuthenticateUser(ctx context.Context, username string, password string) (User, error) {
	getUserFromUsernameQuery := `SELECT *
	                                FROM users
	                                WHERE username = $1 AND users.is_active = true`
	var user User
	err := store.db.QueryRow(ctx, getUserFromUsernameQuery, username).Scan(
		&user.ID,
		&user.Username,
		&user.passwordHash,
		&user.SignUpDate,
		&user.IsActive,
	)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, err
	}

	doesMatch, err := passwords.DoesPasswordMatchHashedPassword(password, user.passwordHash)
	if err != nil {
		return User{}, err
	}

	if doesMatch {
		return user, nil
	} else {
		return User{}, ErrInvalidCredentials
	}
}

// GetUserFromSessionID returns the user attached to the given session id given they have
// a valid session currently in the db.
func (store *UserService) GetUserFromSessionID(ctx context.Context, id string) (User, error) {
	joinSessionAndUserQuery := `SELECT u.id, u.username, u.password_hash, u.sign_up_date, u.is_active
	          FROM sessions s 
	          INNER JOIN users u ON u.id = s.user_id 
	          WHERE s.session_id = $1 AND s.expires_at > NOW() AND u.is_active = true`

	var user User
	err := store.db.QueryRow(ctx, joinSessionAndUserQuery, id).Scan(
		&user.ID,
		&user.Username,
		&user.passwordHash,
		&user.SignUpDate,
		&user.IsActive,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
