package db

import (
	"database/sql"
	"log"
)

func InitDB(db *sql.DB) {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS issues (
		id SERIAL PRIMARY KEY,
		error TEXT,
		cause TEXT,
		fix TEXT
	);
	`)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}
}
