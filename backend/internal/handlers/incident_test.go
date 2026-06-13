package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"devops-memory-assistant/internal/models"
)

// Mock service for testing
type MockIncidentService struct {
	ExtractIncidentFromNotesFunc func(ctx context.Context, rawNotes string) (*models.IncidentInput, error)
	SaveIncidentFunc             func(ctx context.Context, incidentInput *models.IncidentInput) (*models.Incident, error)
	GetIncidentByIDFunc          func(ctx context.Context, id int) (*models.Incident, error)
	UpdateIncidentFunc           func(ctx context.Context, incident *models.Incident) error
	DeleteIncidentFunc           func(ctx context.Context, id int) error
	SearchIncidentsFunc          func(ctx context.Context, query string, limit int) ([]models.Incident, error)
}

func (m *MockIncidentService) ExtractIncidentFromNotes(ctx context.Context, rawNotes string) (*models.IncidentInput, error) {
	if m.ExtractIncidentFromNotesFunc != nil {
		return m.ExtractIncidentFromNotesFunc(ctx, rawNotes)
	}
	return nil, nil
}

func (m *MockIncidentService) SaveIncident(ctx context.Context, incidentInput *models.IncidentInput) (*models.Incident, error) {
	if m.SaveIncidentFunc != nil {
		return m.SaveIncidentFunc(ctx, incidentInput)
	}
	return &models.Incident{}, nil
}

func (m *MockIncidentService) GetIncidentByID(ctx context.Context, id int) (*models.Incident, error) {
	if m.GetIncidentByIDFunc != nil {
		return m.GetIncidentByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockIncidentService) UpdateIncident(ctx context.Context, incident *models.Incident) error {
	if m.UpdateIncidentFunc != nil {
		return m.UpdateIncidentFunc(ctx, incident)
	}
	return nil
}

func (m *MockIncidentService) DeleteIncident(ctx context.Context, id int) error {
	if m.DeleteIncidentFunc != nil {
		return m.DeleteIncidentFunc(ctx, id)
	}
	return nil
}

func (m *MockIncidentService) SearchIncidents(ctx context.Context, query string, limit int) ([]models.Incident, error) {
	if m.SearchIncidentsFunc != nil {
		return m.SearchIncidentsFunc(ctx, query, limit)
	}
	return []models.Incident{}, nil
}

func TestExtractIncident(t *testing.T) {
	// Setup mock service
	mockService := &MockIncidentService{}
	mockService.ExtractIncidentFromNotesFunc = func(ctx context.Context, rawNotes string) (*models.IncidentInput, error) {
		return &models.IncidentInput{
			Title:       "Test Incident",
			Summary:     "Test Summary",
			Symptoms:    []string{"symptom1"},
			Evidence:    []string{"evidence1"},
			RootCause:   []string{"cause1"},
			Resolution:  []string{"fix1"},
			Prevention:  []string{"prevention1"},
			CommandsUsed: []string{"cmd1"},
			Tags:        []string{"tag1"},
			Severity:    "medium",
			Environment: "test",
			ServicesAffected: []string{"service1"},
			LessonsLearned: "lessons learned",
			RawNotes:    rawNotes,
		}, nil
	}

	handler := &IncidentHandler{Service: mockService}

	// Test valid request
	reqBody := bytes.NewBufferString(`{"raw_notes":"Test notes"}`)
	req := httptest.NewRequest(http.MethodPost, "/incident/extract", reqBody)
	w := httptest.NewRecorder()

	handler.ExtractIncident(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var response models.IncidentInput
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Incident", response.Title)
	assert.Equal(t, "Test Summary", response.Summary)
	assert.Equal(t, []string{"symptom1"}, response.Symptoms)
}

func TestExtractIncident_MethodNotAllowed(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incident/extract", nil)
	w := httptest.NewRecorder()

	handler.ExtractIncident(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestExtractIncident_EmptyNotes(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	reqBody := bytes.NewBufferString(`{"raw_notes":""}`)
	req := httptest.NewRequest(http.MethodPost, "/incident/extract", reqBody)
	w := httptest.NewRecorder()

	handler.ExtractIncident(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestExtractIncident_ServiceError(t *testing.T) {
	mockService := &MockIncidentService{}
	mockService.ExtractIncidentFromNotesFunc = func(ctx context.Context, rawNotes string) (*models.IncidentInput, error) {
		return nil, assert.AnError
	}
	handler := &IncidentHandler{Service: mockService}

	reqBody := bytes.NewBufferString(`{"raw_notes":"Test notes"}`)
	req := httptest.NewRequest(http.MethodPost, "/incident/extract", reqBody)
	w := httptest.NewRecorder()

	handler.ExtractIncident(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestSaveIncident(t *testing.T) {
	// Setup mock service
	mockService := &MockIncidentService{}
	mockService.SaveIncidentFunc = func(ctx context.Context, incidentInput *models.IncidentInput) (*models.Incident, error) {
		return &models.Incident{
			ID:          1,
			Title:       incidentInput.Title,
			Summary:     incidentInput.Summary,
			Symptoms:    incidentInput.Symptoms,
			Evidence:    incidentInput.Evidence,
			RootCause:   incidentInput.RootCause,
			Resolution:  incidentInput.Resolution,
			Prevention:  incidentInput.Prevention,
			CommandsUsed: incidentInput.CommandsUsed,
			Tags:        incidentInput.Tags,
			Severity:    incidentInput.Severity,
			Environment: incidentInput.Environment,
			ServicesAffected: incidentInput.ServicesAffected,
			LessonsLearned: incidentInput.LessonsLearned,
			RawNotes:    incidentInput.RawNotes,
		}, nil
	}

	handler := &IncidentHandler{Service: mockService}

	incident := models.Incident{
		Title:       "Test Incident",
		Summary:     "Test Summary",
		Symptoms:    []string{"symptom1"},
		Evidence:    []string{"evidence1"},
		RootCause:   []string{"cause1"},
		Resolution:  []string{"fix1"},
		Prevention:  []string{"prevention1"},
		CommandsUsed: []string{"cmd1"},
		Tags:        []string{"tag1"},
		Severity:    "medium",
		Environment: "test",
		ServicesAffected: []string{"service1"},
		LessonsLearned: "lessons learned",
		RawNotes:    "test notes",
	}

	reqBody, _ := json.Marshal(incident)
	req := httptest.NewRequest(http.MethodPost, "/incident", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	handler.SaveIncident(w, req)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

	var response models.Incident
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test Incident", response.Title)
}

func TestSaveIncident_MethodNotAllowed(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incident", nil)
	w := httptest.NewRecorder()

	handler.SaveIncident(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestSaveIncident_InvalidJSON(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/incident", bytes.NewBufferString(`{invalid json}`))
	w := httptest.NewRecorder()

	handler.SaveIncident(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestGetIncident(t *testing.T) {
	// Setup mock service
	mockService := &MockIncidentService{}
	mockService.GetIncidentByIDFunc = func(ctx context.Context, id int) (*models.Incident, error) {
		return &models.Incident{
			ID:          1,
			Title:       "Test Incident",
			Summary:     "Test Summary",
			Symptoms:    []string{"symptom1"},
			Evidence:    []string{"evidence1"},
			RootCause:   []string{"cause1"},
			Resolution:  []string{"fix1"},
			Prevention:  []string{"prevention1"},
			CommandsUsed: []string{"cmd1"},
			Tags:        []string{"tag1"},
			Severity:    "medium",
			Environment: "test",
			ServicesAffected: []string{"service1"},
			LessonsLearned: "lessons learned",
			RawNotes:    "test notes",
		}, nil
	}

	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incident/1", nil)
	w := httptest.NewRecorder()

	handler.GetIncident(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var response models.Incident
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Test Incident", response.Title)
}

func TestGetIncident_NotFound(t *testing.T) {
	mockService := &MockIncidentService{}
	mockService.GetIncidentByIDFunc = func(ctx context.Context, id int) (*models.Incident, error) {
		return nil, nil
	}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incident/999", nil)
	w := httptest.NewRecorder()

	handler.GetIncident(w, req)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}

func TestGetIncident_InvalidID(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incident/invalid", nil)
	w := httptest.NewRecorder()

	handler.GetIncident(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestUpdateIncident(t *testing.T) {
	// Setup mock service
	mockService := &MockIncidentService{}
	mockService.UpdateIncidentFunc = func(ctx context.Context, incident *models.Incident) error {
		return nil
	}

	handler := &IncidentHandler{Service: mockService}

	incident := models.Incident{
		ID:          1,
		Title:       "Updated Incident",
		Summary:     "Updated Summary",
		Symptoms:    []string{"updated-symptom"},
		Evidence:    []string{"updated-evidence"},
		RootCause:   []string{"updated-cause"},
		Resolution:  []string{"updated-fix"},
		Prevention:  []string{"updated-prevention"},
		CommandsUsed: []string{"updated-cmd"},
		Tags:        []string{"updated-tag"},
		Severity:    "high",
		Environment: "updated",
		ServicesAffected: []string{"updated-service"},
		LessonsLearned: "updated lessons",
		RawNotes:    "updated notes",
	}

	reqBody, _ := json.Marshal(incident)
	req := httptest.NewRequest(http.MethodPut, "/incident/1", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	handler.UpdateIncident(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var response models.Incident
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Updated Incident", response.Title)
}

func TestUpdateIncident_MethodNotAllowed(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/incident/1", nil)
	w := httptest.NewRecorder()

	handler.UpdateIncident(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestUpdateIncident_InvalidJSON(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPut, "/incident/1", bytes.NewBufferString(`{invalid json}`))
	w := httptest.NewRecorder()

	handler.UpdateIncident(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestDeleteIncident(t *testing.T) {
	// Setup mock service
	mockService := &MockIncidentService{}
	mockService.DeleteIncidentFunc = func(ctx context.Context, id int) error {
		return nil
	}

	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/incident/1", nil)
	w := httptest.NewRecorder()

	handler.DeleteIncident(w, req)

	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
}

func TestDeleteIncident_MethodNotAllowed(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incident/1", nil)
	w := httptest.NewRecorder()

	handler.DeleteIncident(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestDeleteIncident_InvalidID(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodDelete, "/incident/invalid", nil)
	w := httptest.NewRecorder()

	handler.DeleteIncident(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestSearchIncidents(t *testing.T) {
	// Setup mock service
	mockService := &MockIncidentService{}
	mockService.SearchIncidentsFunc = func(ctx context.Context, query string, limit int) ([]models.Incident, error) {
		return []models.Incident{
			{
				ID:          1,
				Title:       "Test Incident",
				Summary:     "Test Summary",
				Symptoms:    []string{"symptom1"},
				Evidence:    []string{"evidence1"},
				RootCause:   []string{"cause1"},
				Resolution:  []string{"fix1"},
				Prevention:  []string{"prevention1"},
				CommandsUsed: []string{"cmd1"},
				Tags:        []string{"tag1"},
				Severity:    "medium",
				Environment: "test",
				ServicesAffected: []string{"service1"},
				LessonsLearned: "lessons learned",
				RawNotes:    "test notes",
			},
		}, nil
	}

	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incidents?query=test", nil)
	w := httptest.NewRecorder()

	handler.SearchIncidents(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	var response []models.Incident
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, 1, response[0].ID)
	assert.Equal(t, "Test Incident", response[0].Title)
}

func TestSearchIncidents_MethodNotAllowed(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodPost, "/incidents", nil)
	w := httptest.NewRecorder()

	handler.SearchIncidents(w, req)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
}

func TestSearchIncidents_MissingQuery(t *testing.T) {
	mockService := &MockIncidentService{}
	handler := &IncidentHandler{Service: mockService}

	req := httptest.NewRequest(http.MethodGet, "/incidents", nil)
	w := httptest.NewRecorder()

	handler.SearchIncidents(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}