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
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()

	db.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/issue", handlers.SaveIssue)
	mux.HandleFunc("/search", handlers.SearchIssue)
	mux.HandleFunc("/suggest", handlers.SuggestIssue)

	port := "8080"

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	fmt.Println("Server running on :" + port)

	http.ListenAndServe(":"+port, enableCORS(mux))
}
