package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var version = "v0.2.0"

type Issue struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
	Fix   string `json:"fix"`
}

type Config struct {
	APIURL string `json:"api_url"`
}

// 🔧 Load config
func loadConfig() string {
	home, _ := os.UserHomeDir()
	path := home + "/.devops-memory/config.json"

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("⚠️ Using default localhost API")
		return "http://localhost:8080"
	}
	defer file.Close()

	var config Config
	json.NewDecoder(file).Decode(&config)

	if config.APIURL == "" {
		return "http://localhost:8080"
	}

	return config.APIURL
}

var API = loadConfig()

func main() {

	// 🔥 COMMAND MODE
	if len(os.Args) > 1 {
		command := os.Args[1]

		switch command {

		case "version":
			fmt.Println("DevOps Memory CLI version:", version)
			return

		case "update":
			updateCLI()
			return

		case "search":
			if len(os.Args) < 3 {
				fmt.Println("Usage: devops-memory search <error>")
				return
			}
			runSearch(os.Args[2])
			return

		case "save":
			if len(os.Args) < 5 {
				fmt.Println(`Usage: devops-memory save "error" "cause" "fix"`)
				return
			}
			runSave(os.Args[2], os.Args[3], os.Args[4])
			return
		}
	}

	// 🧠 INTERACTIVE MODE
	runInteractive()
}

func runInteractive() {
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
			fmt.Print("Enter error: ")
			query, _ := reader.ReadString('\n')
			runSearch(strings.TrimSpace(query))

		case "2":
			fmt.Print("Enter error: ")
			errText, _ := reader.ReadString('\n')

			fmt.Print("Enter cause: ")
			causeText, _ := reader.ReadString('\n')

			fmt.Print("Enter fix: ")
			fixText, _ := reader.ReadString('\n')

			runSave(
				strings.TrimSpace(errText),
				strings.TrimSpace(causeText),
				strings.TrimSpace(fixText),
			)

		case "3":
			fmt.Println("Goodbye 👋")
			return

		default:
			color.Red("Invalid choice ❌")
		}
	}
}

// 🔍 SEARCH
func runSearch(query string) {
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

// 💾 SAVE
func runSave(errorText, causeText, fixText string) {

	issue := Issue{
		Error: errorText,
		Cause: causeText,
		Fix:   fixText,
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

	// 🧠 Similar issue detection
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

// 🔄 UPDATE
func updateCLI() {

	fmt.Println("🔄 Updating DevOps Memory CLI...")

	url := "https://github.com/Hardik19102003/devops-memory-assistant/releases/latest/download/devops-memory"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Download error:", err)
		return
	}
	defer resp.Body.Close()

	out, _ := os.Create("/tmp/devops-memory")
	defer out.Close()

	io.Copy(out, resp.Body)

	os.Chmod("/tmp/devops-memory", 0755)

	exec.Command("sudo", "mv", "/tmp/devops-memory", "/usr/local/bin/devops-memory").Run()

	fmt.Println("✅ Updated successfully!")
}