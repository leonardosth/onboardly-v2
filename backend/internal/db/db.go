package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB represents the global database connection pool.
var DB *sql.DB

// InitDB initializes the PostgreSQL connection pool using the pgx driver.
func InitDB(databaseURL string) error {
	var err error
	DB, err = sql.Open("pgx", databaseURL)
	if err != nil {
		return err
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return err
	}

	return nil
}
