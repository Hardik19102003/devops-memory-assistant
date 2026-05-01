package handlers

import (
	"encoding/json"
	"net/http"
)

// ⚡ For now: simple smart suggestions (no AI yet)
func SuggestIssue(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("error")

	// basic intelligence (can replace with AI later)
	suggestions := []string{}

	if query == "CrashLoopBackOff" {
		suggestions = []string{
			"Check pod logs",
			"Verify config/env variables",
			"Check liveness/readiness probes",
		}
	} else if query == "OOMKilled" {
		suggestions = []string{
			"Increase memory limits",
			"Check memory leaks",
			"Optimize application usage",
		}
	}

	json.NewEncoder(w).Encode(suggestions)
}