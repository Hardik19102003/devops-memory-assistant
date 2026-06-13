package main

import (
	"fmt"
	"log"
	"net/http"

	"devops-memory-assistant/config"
	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/handlers"
	"devops-memory-assistant/internal/service"

	"github.com/joho/godotenv"
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
    err := godotenv.Load()
    if err != nil {
        log.Println("No .env file found, using system env")
    }
	cfg := config.Load()

	fmt.Println("DB:", cfg.DBURL)
	fmt.Println("PORT:", cfg.PORT)

	// 🔥 PASS CONFIG TO DB (important change)
	database := db.New(cfg.DBURL)

	db.RunMigrations(database.DB)

	// Initialize services
	incidentService := service.NewIncidentService()
	incidentHandler := handlers.NewIncidentHandler(incidentService)

	mux := http.NewServeMux()
	// Old endpoints (keep for backward compatibility or remove if not needed)
	mux.HandleFunc("/issue", handlers.SaveIssue)
	mux.HandleFunc("/search", handlers.SearchIssue)
	mux.HandleFunc("/suggest", handlers.SuggestIssue)
	mux.HandleFunc("/delete", handlers.DeleteIssue)
	mux.HandleFunc("/analyze", handlers.AnalyzeIssue)
	
	// New incident endpoints
	mux.HandleFunc("/incident/extract", incidentHandler.ExtractIncident)
	mux.HandleFunc("/incident", incidentHandler.SaveIncident)
	mux.HandleFunc("/incident/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			incidentHandler.GetIncident(w, r)
		case http.MethodPut:
			incidentHandler.UpdateIncident(w, r)
		case http.MethodDelete:
			incidentHandler.DeleteIncident(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/incidents", incidentHandler.SearchIncidents)

	fmt.Println("Server running on :", cfg.PORT)

	err = http.ListenAndServe(":"+cfg.PORT, enableCORS(mux))
	if err != nil {
		panic(err)
	}
}