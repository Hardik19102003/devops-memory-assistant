package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/fatih/color"
)

type Issue struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
	Fix   string `json:"fix"`
}

func main() {

	// ⚙️ FLAGS
	limit := flag.Int("limit", 5, "number of results")
	jsonOutput := flag.Bool("json", false, "output in JSON format")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage:")
		fmt.Println("  devops-memory search <error>")
		fmt.Println("  devops-memory save <error> <cause> <fix>")
		return
	}

	command := args[0]

	// 🎨 COLORS
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	API := "http://localhost:8080" // change to deployed later

	// =========================
	// 🔍 SEARCH COMMAND
	// =========================
	if command == "search" {

		if len(args) < 2 {
			fmt.Println("Usage: devops-memory search <error>")
			return
		}

		query := args[1]

		url := fmt.Sprintf("%s/search?error=%s", API, query)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		var issues []Issue
		err = json.NewDecoder(resp.Body).Decode(&issues)
		if err != nil {
			fmt.Println("Failed to parse response")
			return
		}

		if len(issues) > *limit {
			issues = issues[:*limit]
		}

		if *jsonOutput {
			json.NewEncoder(os.Stdout).Encode(issues)
			return
		}

		fmt.Println(magenta("🔍 DevOps Memory Assistant Results\n"))

		if len(issues) == 0 {
			color.Yellow("No results found 👀")
			return
		}

		for _, issue := range issues {
			fmt.Println("--------------------------------------------------")
			fmt.Println(cyan("Error:"), red(issue.Error))
			fmt.Println(cyan("Cause:"), issue.Cause)
			fmt.Println(cyan("Fix:"), green(issue.Fix))
		}
	}

	// =========================
	// 💾 SAVE COMMAND
	// =========================
	if command == "save" {

		if len(args) < 4 {
			fmt.Println(`Usage: devops-memory save "error" "cause" "fix"`)
			return
		}

		errorText := args[1]
		causeText := args[2]
		fixText := args[3]

		issue := Issue{
			Error: errorText,
			Cause: causeText,
			Fix:   fixText,
		}

		body, _ := json.Marshal(issue)

		resp, err := http.Post(API+"/issue", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		// 🧠 HANDLE SIMILAR ISSUE
		if existing, ok := response["existing"]; ok {
			fmt.Println(color.YellowString("⚠️ Similar issue already exists!\n"))

			ex := existing.(map[string]interface{})

			fmt.Println(cyan("Error:"), red(ex["error"]))
			fmt.Println(cyan("Cause:"), ex["cause"])
			fmt.Println(cyan("Fix:"), green(ex["fix"]))
			return
		}

		// ✅ SUCCESS
		fmt.Println(color.GreenString("✅ Issue saved successfully!"))
	}
}