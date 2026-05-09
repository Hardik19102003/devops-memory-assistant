package main

import (
	"fmt"
	"net/http"
	"os"

	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/handlers"

	"github.com/joho/godotenv"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)

		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, DELETE, OPTIONS",
		)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()

	database := db.New()

	db.RunMigrations(database.DB)

	mux := http.NewServeMux()
	mux.HandleFunc("/issue", handlers.SaveIssue)
	mux.HandleFunc("/search", handlers.SearchIssue)
	mux.HandleFunc("/suggest", handlers.SuggestIssue)
	mux.HandleFunc("/delete", handlers.DeleteIssue)

	port := "8080"

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Println("Server running on :" + port)

	err := http.ListenAndServe(":"+port, enableCORS(mux))
if err != nil {
	panic(err)
}
}
