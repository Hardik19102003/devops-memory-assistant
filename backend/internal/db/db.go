package db

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
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

func RunMigrations(db *sql.DB) {
	err := goose.Up(db, "../../migrations")
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
}
