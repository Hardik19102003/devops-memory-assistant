package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"devops-memory-assistant/internal/models"
)

// ExtractIncidentFromNotes extracts structured incident data from raw notes using Ollama.
// cleanJSONResponse strips markdown code fences and extracts the first JSON object/array.
func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)
	// Strip markdown code fences if present
	if strings.HasPrefix(response, "```") {
		lines := strings.Split(response, "\n")
		var buf strings.Builder
		for i, line := range lines {
			if i == 0 && strings.HasPrefix(line, "```") {
				continue
			}
			if strings.TrimSpace(line) == "```" {
				continue
			}
			buf.WriteString(line)
			buf.WriteByte('\n')
		}
		response = strings.TrimSpace(buf.String())
	}
	return response
}

func normalizeStringArray(raw []byte) []string {
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr
	}
	// Try array of objects → stringify each object
	var objArr []map[string]interface{}
	if err := json.Unmarshal(raw, &objArr); err == nil {
		for _, obj := range objArr {
			b, _ := json.Marshal(obj)
			arr = append(arr, string(b))
		}
		return arr
	}
	// Try single object → stringify
	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err == nil {
		b, _ := json.Marshal(obj)
		return []string{string(b)}
	}
	// Try single string
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return []string{s}
	}
	return nil
}

func parseIncidentInput(data []byte) (*models.IncidentInput, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	var in models.IncidentInput
	json.Unmarshal(raw["title"], &in.Title)
	json.Unmarshal(raw["summary"], &in.Summary)
	in.Symptoms = normalizeStringArray(raw["symptoms"])
	in.Evidence = normalizeStringArray(raw["evidence"])
	in.RootCause = normalizeStringArray(raw["root_cause"])
	in.Resolution = normalizeStringArray(raw["resolution"])
	in.Prevention = normalizeStringArray(raw["prevention"])
	in.CommandsUsed = normalizeStringArray(raw["commands_used"])
	in.Tags = normalizeStringArray(raw["tags"])
	json.Unmarshal(raw["severity"], &in.Severity)
	json.Unmarshal(raw["environment"], &in.Environment)
	in.ServicesAffected = normalizeStringArray(raw["services_affected"])
	json.Unmarshal(raw["lessons_learned"], &in.LessonsLearned)
	json.Unmarshal(raw["raw_notes"], &in.RawNotes)
	return &in, nil
}

func ExtractIncidentFromNotes(notes string) (*models.IncidentInput, error) {
	prompt := `
You are a DevOps incident analyst. Extract the following fields from the incident notes and return them as a JSON object:
- title: a short title summarizing the incident
- summary: a brief summary of the incident
- symptoms: array of plain strings
- evidence: array of plain strings
- root_cause: array of plain strings
- resolution: array of plain strings
- prevention: array of plain strings
- commands_used: array of plain strings
- tags: array of plain strings
- severity: one of low, medium, high, critical
- environment: string
- services_affected: array of plain strings
- lessons_learned: string

Return ONLY a JSON object with these fields. All array fields must be arrays of plain strings, not objects. If a field is not present, return an empty array or empty string as appropriate.

Incident notes: ` + notes

	reqBody := OllamaGenerateRequest{
		Model:       "phi3:mini",
		Prompt:      prompt,
		Stream:      false,
		Temperature: 0.1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call Ollama: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result OllamaGenerateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Ollama response: %w", err)
	}

	// The response should be a JSON string. Let's unmarshal it into IncidentInput.
	cleaned := cleanJSONResponse(result.Response)
	incidentInput, err := parseIncidentInput([]byte(cleaned))
	if err != nil {
		return nil, fmt.Errorf("failed to parse incident data: %w", err)
	}

	return incidentInput, nil
}