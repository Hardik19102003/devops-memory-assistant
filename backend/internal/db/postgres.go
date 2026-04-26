package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

var DB *sql.DB

func New() *Database {
	connStr := "user=devops password=devops dbname=devops_memory sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Connected to PostgreSQL ✅")

	DB = db
	return &Database{DB: db}
}
