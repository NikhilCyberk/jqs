package models

import (
	"database/sql"
)

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
		id SERIAL PRIMARY KEY,
		payload JSONB NOT NULL,
		status VARCHAR(20) NOT NULL,
		result JSONB,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`
	_, err := db.Exec(query)
	return err
}
