package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL not set")
	}

	db, err := sql.Open("pgx", dbUrl)

	if err != nil {
		return nil, fmt.Errorf("Coulnt connect to postgres %v", err)
	}

	// Check if DB is reachable
	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("Unreachable %v", err)
	}
	return db, nil

}
