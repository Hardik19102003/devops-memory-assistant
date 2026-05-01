package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

var DB *sql.DB

func New() *Database {
	log.Println("DB URL:", os.Getenv("DATABASE_URL"))
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=devops password=devops dbname=devops_memory sslmode=disable host=localhost"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("Connected to PostgreSQL ✅")

	// InitDB(db) // 👈 REMOVE THIS LINE - migrations handle it now

	DB = db
	return &Database{DB: db}
}
