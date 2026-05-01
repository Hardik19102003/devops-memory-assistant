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

	if issue.Error == "" || issue.Cause == "" || issue.Fix == "" {
		http.Error(w, "All fields required", 400)
		return
	}

	// 🔍 Step 1: Check similar issue
	existing, err := db.FindSimilarIssue(issue.Error)

	if err == nil && existing != nil {
		// 👇 Return suggestion instead of saving
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Similar issue found",
			"existing": existing,
		})
		return
	}

	// 💾 Step 2: Save if not found
	err = db.SaveIssue(issue)

	if err != nil {
		http.Error(w, "Error saving", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Saved successfully",
	})
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
