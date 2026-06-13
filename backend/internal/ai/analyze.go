package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func AnalyzeIncident(query string, context string) (string, error) {

	prompt := `
You are a senior Kubernetes and DevOps incident investigator.

IMPORTANT RULES:
- ONLY use evidence from timeline
- NEVER invent information
- NEVER mention technologies not present
- NEVER generate stories
- NEVER generate unrelated text
- Keep answers SHORT and technical
- If evidence missing say:
NOT ENOUGH EVIDENCE

Return EXACTLY in this format:

ROOT CAUSE:
- ...

EVIDENCE:
- ...

FIX:
- ...

PREVENTION:
- ...

TIMELINE:
` + context + `

USER QUERY:
` + query

	reqBody := OllamaGenerateRequest{
		Model:       "phi3:mini",
		Prompt:      prompt,
		Stream:      false,
		Temperature: 0.1,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result OllamaGenerateResponse

	err = json.Unmarshal(body, &result)

	if err != nil {
		return "", err
	}

	return result.Response, nil
}