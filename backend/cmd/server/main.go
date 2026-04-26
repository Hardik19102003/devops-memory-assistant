package main

import (
	"fmt"
	"net/http"

	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/handlers"
)

func main() {

	db.New()

	http.HandleFunc("/issue", handlers.SaveIssue)
	http.HandleFunc("/search", handlers.SearchIssue)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
