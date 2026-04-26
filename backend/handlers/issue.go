package handlers

import (
	"encoding/json"
	"net/http"
)


type Issue struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
	Fix   string `json:"fix"`
}

var issues  []Issue

func SaveIssue(w http.ResponseWriter, r *http.Request) {
	var issue Issue
	json.NewDecoder(r.Body).Decode(&issue)
	issues = append(issues, issue)

	json.NewEncoder(w).Encode(map[string]string{"message": "Issue saved successfully",})}
