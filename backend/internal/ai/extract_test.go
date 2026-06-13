package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"devops-memory-assistant/internal/models"
)

// Helper to create test incident input
func createTestIncidentInput() *models.IncidentInput {
	return &models.IncidentInput{
		Title:       "CrashLoopBackOff after Secret Rotation",
		Summary:     "Pods entered CrashLoopBackOff state following a secret rotation event due to accidental secret deletion.",
		Symptoms:    []string{"Pods in CrashLoopBackOff", "Secret mount failures"},
		Evidence:    []string{"kubectl describe pod showed secret mount failures"},
		RootCause:   []string{"Secret was accidentally deleted during rotation"},
		Resolution:  []string{"Recreated the secret", "Restarted the deployment"},
		Prevention:  []string{"Implement secret rotation verification steps", "Use secrets management tool with rollback capability"},
		CommandsUsed: []string{"kubectl describe pod", "kubectl create secret", "kubectl rollout restart deployment"},
		Tags:        []string{"kubernetes", "secret", "crashloop"},
		Severity:    "medium",
		Environment: "production",
		ServicesAffected: []string{"web-api", "auth-service"},
		LessonsLearned: "Always verify secret existence post-rotation before considering rotation complete",
		RawNotes:    "Pods entered CrashLoopBackOff after secret rotation. kubectl describe pod showed secret mount failures. Secret was accidentally deleted. Recreated secret and restarted deployment. Service recovered.",
	}
}

func TestExtractIncidentFromNotes(t *testing.T) {
	// Setup mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock Ollama response
	httpmock.RegisterResponder("POST", "http://localhost:11434/api/generate",
		func(req *http.Request) (*http.Response, error) {
			testInput := createTestIncidentInput()
			jsonBytes, _ := json.Marshal(testInput)
			response := string(`{"response": ` + string(jsonBytes) + `}`)
			resp := httptest.NewRecorder()
			resp.Body.Write([]byte(response))
			resp.Header().Set("Content-Type", "application/json")
			return resp.Result(), nil
		})

	// Test input
	notes := "Pods entered CrashLoopBackOff after secret rotation. kubectl describe pod showed secret mount failures. Secret was accidentally deleted. Recreated secret and restarted deployment. Service recovered."

	// Call function
	result, err := ExtractIncidentFromNotes(notes)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "CrashLoopBackOff after Secret Rotation", result.Title)
	assert.Equal(t, "Pods entered CrashLoopBackOff state following a secret rotation event due to accidental secret deletion.", result.Summary)
	assert.Equal(t, []string{"Pods in CrashLoopBackOff", "Secret mount failures"}, result.Symptoms)
	assert.Equal(t, []string{"kubectl describe pod showed secret mount failures"}, result.Evidence)
	assert.Equal(t, []string{"Secret was accidentally deleted during rotation"}, result.RootCause)
	assert.Equal(t, []string{"Recreated the secret", "Restarted the deployment"}, result.Resolution)
	assert.Equal(t, []string{"Implement secret rotation verification steps", "Use secrets management tool with rollback capability"}, result.Prevention)
	assert.Equal(t, []string{"kubectl describe pod", "kubectl create secret", "kubectl rollout restart deployment"}, result.CommandsUsed)
	assert.Equal(t, []string{"kubernetes", "secret", "crashloop"}, result.Tags)
	assert.Equal(t, "medium", result.Severity)
	assert.Equal(t, "production", result.Environment)
	assert.Equal(t, []string{"web-api", "auth-service"}, result.ServicesAffected)
	assert.Equal(t, "Always verify secret existence post-rotation before considering rotation complete", result.LessonsLearned)
	assert.Equal(t, notes, result.RawNotes)
}

func TestExtractIncidentFromNotes_InvalidJSON(t *testing.T) {
	// Setup mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock Ollama response with invalid JSON
	httpmock.RegisterResponder("POST", "http://localhost:11434/api/generate",
		func(req *http.Request) (*http.Response, error) {
			response := `{"response": "invalid json"}`
			resp := httptest.NewRecorder()
			resp.Body.Write([]byte(response))
			resp.Header().Set("Content-Type", "application/json")
			return resp.Result(), nil
		})

	// Test input
	notes := "Some notes"

	// Call function
	_, err := ExtractIncidentFromNotes(notes)

	// Assert error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal incident data")
}

func TestExtractIncidentFromNotes_OllamaError(t *testing.T) {
	// Setup mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock Ollama service error
	httpmock.RegisterResponder("POST", "http://localhost:11434/api/generate",
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       http.NoBody,
				Header:     make(http.Header),
			}, nil
		})

	// Test input
	notes := "Some notes"

	// Call function
	_, err := ExtractIncidentFromNotes(notes)

	// Assert error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to call Ollama")
}