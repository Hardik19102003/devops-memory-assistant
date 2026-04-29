package handlers

import (
	"encoding/json"
	"net/http"

	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/models"
)

func SaveIssue(w http.ResponseWriter, r *http.Request) {

	var issue models.Issue

	json.NewDecoder(r.Body).Decode(&issue)

	err := db.SaveIssue(issue)

	if err != nil {
		http.Error(w, "Error saving", 500)
		return
	}

	json.NewEncoder(w).Encode("Saved successfully")
}

func SearchIssue(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("error")

	if query == "" {
		http.Error(w, "Missing search query", 400)
		return
	}

	results, err := db.SearchIssue(query)

	if err != nil {
		http.Error(w, "Error fetching data", 500)
		return
	}
	if results == nil {
    	results = []models.Issue{}
	}
	json.NewEncoder(w).Encode(results)
}