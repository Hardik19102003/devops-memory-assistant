package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/models"
)

func SaveIssue(w http.ResponseWriter, r *http.Request) {

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

	if issue.Error == "" {
	log.Printf("WARN: missing error field")
	http.Error(w, "Error field required", 400)
	return
}

if len(issue.Causes) == 0 {
	http.Error(w, "At least one cause required", 400)
	return
}

if len(issue.Fixes) == 0 {
	http.Error(w, "At least one fix required", 400)
	return
}

	// 🔍 Check similar issue
	existing, err := db.FindSimilarIssue(issue.Error)

	similarFound := false

	if err == nil && existing != nil {
		log.Printf(
			"INFO: similar issue found for error: %s",
			issue.Error,
		)

		similarFound = true
	}

	// 💾 Save anyway
	err = db.SaveIssue(issue)

	if err != nil {
		log.Printf("ERROR: failed to save issue: %v", err)
		http.Error(w, "Error saving", 500)
		return
	}

	log.Printf("INFO: successfully saved issue: %s", issue.Error)

	// ✅ Single response only
	if similarFound {

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Saved successfully, but similar issue exists ⚠️",
			"similar": existing,
		})

		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Saved successfully ✅",
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

func DeleteIssue(w http.ResponseWriter, r *http.Request) {

	apiKey := r.Header.Get("Authorization")

	if apiKey != "Bearer devops-secret-key" {
		http.Error(w, "Unauthorized", 401)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing issue ID", 400)
		return
	}

	err := db.DeleteIssue(id)

	if err != nil {
		log.Printf("ERROR: failed to delete issue: %v", err)

		http.Error(w, "Delete failed", 500)
		return
	}

	log.Printf("INFO: issue deleted successfully: %s", id)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Issue deleted successfully ✅",
	})
}
