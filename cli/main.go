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
	ID        int      `json:"id"`
	Error     string   `json:"error"`
	Cause     string   `json:"cause"`
	Fix       string   `json:"fix"`
	Steps     string   `json:"steps"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
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

func getAPI() string {
	return loadConfig()
}

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
	if len(os.Args) < 7 {
		fmt.Println(`Usage:
devops-memory save "error" "cause" "fix" "steps" "tag1,tag2"
`)
		return

	}

	tags := strings.Split(os.Args[6], ",")

	runSave(
		os.Args[2],
		os.Args[3],
		os.Args[4],
		os.Args[5],
		tags,
	)

	return

	case "delete":
	if len(os.Args) < 3 {
		fmt.Println("Usage: devops-memory delete <id>")
		return
	}

	id := os.Args[2]
	runDelete(id)
	return
		}
	}

	// 🧠 INTERACTIVE MODE
	runInteractive()
}

func runDelete(id string) {

	url := fmt.Sprintf("%s/delete?id=%s", getAPI(), id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer devops-secret-key")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))

	if resp.StatusCode == 200 {
		color.Green("✅ Issue deleted successfully")
	} else {
		color.Red("❌ Failed to delete issue")
	}
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

fmt.Print("Enter steps: ")
stepsText, _ := reader.ReadString('\n')

fmt.Print("Enter tags (comma separated): ")
tagsText, _ := reader.ReadString('\n')

tags := strings.Split(strings.TrimSpace(tagsText), ",")

runSave(
	strings.TrimSpace(errText),
	strings.TrimSpace(causeText),
	strings.TrimSpace(fixText),
	strings.TrimSpace(stepsText),
	tags,
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
	url := fmt.Sprintf("%s/search?error=%s", getAPI(), query)

	resp, err := http.Get(url)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	var issues []Issue
	err = json.NewDecoder(resp.Body).Decode(&issues)
if err != nil {
	color.Red("JSON Decode Error: %v", err)
	return
}

fmt.Printf("%+v\n", issues)

	if len(issues) == 0 {
		color.Yellow("No results found 👀")
	} else {
		for _, issue := range issues {

			if issue.Cause == "" || issue.Fix == "" {
				continue // skip bad data
			}

			fmt.Println("\n---------------------------")

color.Cyan("🚨 Error: %s", issue.Error)

fmt.Println("📌 Cause:")
fmt.Println(issue.Cause)

color.Green("\n✅ Fix:")
fmt.Println(issue.Fix)

if issue.Steps != "" {
	color.Yellow("\n🛠 Steps:")
	fmt.Println(issue.Steps)
}

if len(issue.Tags) > 0 {
	color.Magenta("\n🏷 Tags:")
	fmt.Println(strings.Join(issue.Tags, ", "))
}

fmt.Println("\n🕒 Created:", issue.CreatedAt)
		}
	}

	// 🔥 NEW: Fetch suggestions from backend
	suggestURL := fmt.Sprintf("%s/suggest?error=%s", getAPI(), query)

	resp2, err := http.Get(suggestURL)
	if err != nil {
		return // silently skip if suggestion fails
	}
	defer resp2.Body.Close()

	var suggestions []string
	json.NewDecoder(resp2.Body).Decode(&suggestions)

	if len(suggestions) > 0 {
		color.Magenta("\n💡 Suggestions:")
		for _, s := range suggestions {
			fmt.Println("👉", s)
		}
	}
}

// 💾 SAVE
func runSave(
	errorText string,
	causeText string,
	fixText string,
	stepsText string,
	tags []string,
) {

	issue := Issue{
	Error: errorText,
	Cause: causeText,
	Fix:   fixText,
	Steps: stepsText,
	Tags:  tags,
}

	body, _ := json.Marshal(issue)

	resp, err := http.Post(getAPI()+"/issue", "application/json", bytes.NewBuffer(body))
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
