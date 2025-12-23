// Package store provides functionality for CRUD on all of the tables in DB.
package store

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	db *pgxpool.Pool
}

func (store *UserStore) DoesUserWithUsernameExist(username string, c *gin.Context) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)"

	var doesUsernameAlreadyExist bool
	err := store.db.QueryRow(c.Request.Context(), query, username).Scan(&doesUsernameAlreadyExist)
	if err != nil {
		return false, err
	}

	return doesUsernameAlreadyExist, nil
}
