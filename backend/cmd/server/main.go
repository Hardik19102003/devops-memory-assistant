package main

import (
	"fmt"
	"net/http"

	"devops-memory-assistant/config"
	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/handlers"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Load()

	fmt.Println("DB:", cfg.DBURL)
	fmt.Println("PORT:", cfg.PORT)

	// 🔥 PASS CONFIG TO DB (important change)
	database := db.New(cfg.DBURL)

	db.RunMigrations(database.DB)

	mux := http.NewServeMux()
	mux.HandleFunc("/issue", handlers.SaveIssue)
	mux.HandleFunc("/search", handlers.SearchIssue)
	mux.HandleFunc("/suggest", handlers.SuggestIssue)
	mux.HandleFunc("/delete", handlers.DeleteIssue)

	fmt.Println("Server running on :", cfg.PORT)

	err := http.ListenAndServe(":"+cfg.PORT, enableCORS(mux))
	if err != nil {
		panic(err)
	}
}
