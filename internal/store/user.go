// Package store provides the service layer exposing databse operations.
package store

import (
	"context"

	"gochat/main/internal/models"
	"gochat/main/internal/utils/passwords"

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

// CreateUser returns true if password hash matches the hash stored in the DB.
func (store *UserService) CreateUser(username string, password string, context context.Context) (models.User, error) {
	passHash, err := passwords.CreatePasswordHash(password, passwords.DefaultArgon2Params)
	if err != nil {
		return models.User{}, err
	}

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING *"

	var user models.User
	err = store.db.QueryRow(context, query, username, passHash).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.SignUpDate,
		&user.IsActive,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// AuthenticateUser returns true if the user can be authenticated with the given
// credentials, otherwise false and an optional error.
func (store *UserService) AuthenticateUser(username string, password string, context context.Context) (bool, error) {
	getPasswordHashQuery := "SELECT password_hash FROM users WHERE username = $1"
	var passwordHash string
	err := store.db.QueryRow(context, getPasswordHashQuery, username).Scan(&passwordHash)
	if err != nil {
		return false, err
	}

	doesMatch, err := passwords.DoesPasswordMatchHashedPassword(password, passwordHash)
	if err != nil {
		return false, err
	}

	return doesMatch, nil
}
