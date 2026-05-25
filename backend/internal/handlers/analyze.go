package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"devops-memory-assistant/internal/ai"
	"devops-memory-assistant/internal/db"
)

func AnalyzeIssue(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("query")

	if query == "" {

		http.Error(w, "missing query", http.StatusBadRequest)

		return
	}

	log.Println("STEP 1: finding related incidents")

	relatedIncidents, err := db.FindRelatedIncidents(query)

	if err != nil {

		log.Println("ERROR:", err)

		http.Error(w, "failed retrieval", http.StatusInternalServerError)

		return
	}

	if len(relatedIncidents) == 0 {

		http.Error(w, "no related incidents found", http.StatusNotFound)

		return
	}

	// Best semantic match
	bestMatch := relatedIncidents[0]

	// Confidence scoring
	confidence := ai.CalculateConfidence(bestMatch.Distance)

	// Extract only issue IDs
	var issueIDs []int

	for _, incident := range relatedIncidents {

		issueIDs = append(issueIDs, incident.IssueID)
	}

	log.Println("STEP 2: building multi-incident context")

	context, err := db.BuildMultiIncidentContext(issueIDs)

	if err != nil {

		log.Println("ERROR:", err)

		http.Error(w, "failed context build", http.StatusInternalServerError)

		return
	}

	log.Println("STEP 3: calling AI")

	analysis, err := ai.AnalyzeIncident(query, context)

	if err != nil {

		log.Println("ERROR:", err)

		http.Error(w, "failed AI analysis", http.StatusInternalServerError)

		return
	}

	log.Println("STEP 4: AI response received")

	response := map[string]interface{}{
		"query": query,

		"matched_incident": bestMatch.Error,

		"matched_issue_id": bestMatch.IssueID,

		"similarity_score": bestMatch.Distance,

		"confidence": confidence,

		"related_incidents": relatedIncidents,

		"context": context,

		"analysis": analysis,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}