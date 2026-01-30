package db

import "database/sql"

func RunMigrations(db *sql.DB) error {
	userTable := `
	CREATE TABLE IF NOT EXISTS users(
	  id SERIAL PRIMARY KEY,
	  name VARCHAR(20) NOT NULL,
	  email VARCHAR(50) NOT NULL UNIQUE,
	  password VARCHAR(100) NOT NULL,
	  created_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Exec(userTable)
	return err

}
