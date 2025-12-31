// Package store provides functionality for CRUD on all of the tables in DB.
package store

import (
	"gochat/main/internal/models"
	"gochat/main/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserService {
	return &UserService{
		db: db,
	}
}

// CreateUser returns true if password hash matches the hash stored in the DB.
func (store *UserService) CreateUser(username string, password string, c *gin.Context) (models.User, error) {
	passHash, err := utils.CreatePasswordHash(password, utils.DefaultArgon2Params)
	if err != nil {
		return models.User{}, err
	}

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING *"

	var user models.User
	err = store.db.QueryRow(c.Request.Context(), query, username, passHash).Scan(
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
