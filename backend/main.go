package main

import (
	"devops-memory-assistant/backend/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "DevOps Memory Assistant running 🚀")
	})

	http.HandleFunc("/issue", handlers.SaveIssue)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}