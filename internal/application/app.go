package application

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// App stores the resorces shared between requests.
type App struct {
	DB *pgxpool.Pool
}
