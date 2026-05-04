package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/models"
)

func SaveIssue(w http.ResponseWriter, r *http.Request) {
	// Simple API key auth
	apiKey := r.Header.Get("Authorization")
	if apiKey != "Bearer devops-secret-key" {
		log.Printf("WARN: unauthorized access attempt")
		http.Error(w, "Unauthorized", 401)
		return
	}

	var issue models.Issue

	if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
		log.Printf("ERROR: failed to decode request body: %v", err)
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if issue.Error == "" || issue.Cause == "" || issue.Fix == "" {
		log.Printf("WARN: missing required fields in issue save attempt")
		http.Error(w, "All fields required", 400)
		return
	}

	// 🔍 Step 1: Check similar issue
	existing, err := db.FindSimilarIssue(issue.Error)

	if err == nil && existing != nil {
		log.Printf("INFO: similar issue found for error: %s", issue.Error)
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
		log.Printf("ERROR: failed to save issue: %v", err)
		http.Error(w, "Error saving", 500)
		return
	}

	log.Printf("INFO: successfully saved issue: %s", issue.Error)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Saved successfully",
	})
}

func SearchIssue(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("error")

	if query == "" {
		log.Printf("WARN: missing search query parameter")
		http.Error(w, "Missing search query", 400)
		return
	}

	results, err := db.SearchIssue(query)

	if err != nil {
		log.Printf("ERROR: failed to search issues: %v", err)
		http.Error(w, "Error fetching data", 500)
		return
	}
	if results == nil {
		results = []models.Issue{}
	}

	log.Printf("INFO: search completed for query '%s', found %d results", query, len(results))
	json.NewEncoder(w).Encode(results)
}
