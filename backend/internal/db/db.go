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
	log.Println("DB connected (URL loaded)")
}

func RunMigrations(db *sql.DB) {

	paths := []string{
		"../../migrations", // local dev
		"./migrations",     // production/render
		"migrations",       // fallback
	}

	var err error

	for _, path := range paths {

		err = goose.Up(db, path)

		if err == nil {
			log.Println("Migrations ran successfully using:", path)
			return
		}

		log.Println("Migration path failed:", path)
	}

	log.Fatal("Migration failed:", err)
}
