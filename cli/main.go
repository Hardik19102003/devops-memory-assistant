package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Issue struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
	Fix   string `json:"fix"`
}

type Config struct {
	APIURL string `json:"api_url"`
}

func loadConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "http://localhost:8080"
	}

	path := home + "/.devops-memory/config.json"

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("⚠️ Config not found, using default localhost")
		return "http://localhost:8080"
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println("⚠️ Invalid config, using default localhost")
		return "http://localhost:8080"
	}

	return config.APIURL
}

var API = loadConfig()

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		color.Magenta("\n🚀 DevOps Memory Assistant\n")
		fmt.Println("1. Search Issue")
		fmt.Println("2. Save Issue")
		fmt.Println("3. Exit")

		fmt.Print("\nEnter choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {

		case "1":
			handleSearch(reader)

		case "2":
			handleSave(reader)

		case "3":
			fmt.Println("Goodbye 👋")
			return

		default:
			color.Red("Invalid choice ❌")
		}
	}
}

func handleSearch(reader *bufio.Reader) {

	fmt.Print("Enter error: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	url := fmt.Sprintf("%s/search?error=%s", API, query)

	resp, err := http.Get(url)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	var issues []Issue
	json.NewDecoder(resp.Body).Decode(&issues)

	if len(issues) == 0 {
		color.Yellow("No results found 👀")
		return
	}

	for _, issue := range issues {
		fmt.Println("\n---------------------------")
		color.Cyan("Error: %s", issue.Error)
		fmt.Println("Cause:", issue.Cause)
		color.Green("Fix: %s", issue.Fix)
	}
}

func handleSave(reader *bufio.Reader) {

	fmt.Print("Enter error: ")
	errorText, _ := reader.ReadString('\n')

	fmt.Print("Enter cause: ")
	causeText, _ := reader.ReadString('\n')

	fmt.Print("Enter fix: ")
	fixText, _ := reader.ReadString('\n')

	issue := Issue{
		Error: strings.TrimSpace(errorText),
		Cause: strings.TrimSpace(causeText),
		Fix:   strings.TrimSpace(fixText),
	}

	body, _ := json.Marshal(issue)

	resp, err := http.Post(API+"/issue", "application/json", bytes.NewBuffer(body))
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if existing, ok := response["existing"]; ok {
		color.Yellow("\n⚠️ Similar issue already exists!\n")

		ex := existing.(map[string]interface{})

		color.Cyan("Error: %v", ex["error"])
		fmt.Println("Cause:", ex["cause"])
		color.Green("Fix: %v", ex["fix"])
		return
	}

	color.Green("\n✅ Issue saved successfully!")
}