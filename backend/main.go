package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "DevOps Memory Assistant 🚀")
	})

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}