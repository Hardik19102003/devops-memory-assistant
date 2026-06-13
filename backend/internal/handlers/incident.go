package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"devops-memory-assistant/internal/models"
)

// IncidentService defines the interface for incident business logic.
type IncidentService interface {
	ExtractIncidentFromNotes(ctx context.Context, rawNotes string) (*models.IncidentInput, error)
	SaveIncident(ctx context.Context, incidentInput *models.IncidentInput) (*models.Incident, error)
	GetIncidentByID(ctx context.Context, id int) (*models.Incident, error)
	UpdateIncident(ctx context.Context, incident *models.Incident) error
	DeleteIncident(ctx context.Context, id int) error
	SearchIncidents(ctx context.Context, query string, limit int) ([]models.Incident, error)
}

// IncidentHandler handles HTTP requests for incidents.
type IncidentHandler struct {
	Service IncidentService
}

// NewIncidentHandler creates a new IncidentHandler.
func NewIncidentHandler(service IncidentService) *IncidentHandler {
	return &IncidentHandler{Service: service}
}

// ExtractIncident handles POST /incident/extract
// It takes raw notes and returns a structured incident for preview.
func (h *IncidentHandler) ExtractIncident(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RawNotes string `json:"raw_notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.RawNotes == "" {
		http.Error(w, "raw_notes cannot be empty", http.StatusBadRequest)
		return
	}

	incidentInput, err := h.Service.ExtractIncidentFromNotes(r.Context(), req.RawNotes)
	if err != nil {
		log.Printf("failed to extract incident: %v", err)
		http.Error(w, "failed to extract incident", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incidentInput)
}

// SaveIncident handles POST /incidents
// It saves a structured incident (after user edits).
func (h *IncidentHandler) SaveIncident(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var incident models.IncidentInput
	if err := json.NewDecoder(r.Body).Decode(&incident); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	savedIncident, err := h.Service.SaveIncident(r.Context(), &incident)
	if err != nil {
		log.Printf("failed to save incident: %v", err)
		http.Error(w, "failed to save incident", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(savedIncident)
}

// GetIncident handles GET /incidents/{id}
func (h *IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/incident/"):]
	if idStr == "" {
		http.Error(w, "missing incident id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid incident id", http.StatusBadRequest)
		return
	}

	incident, err := h.Service.GetIncidentByID(r.Context(), id)
	if err != nil {
		log.Printf("failed to get incident: %v", err)
		http.Error(w, "failed to get incident", http.StatusInternalServerError)
		return
	}
	if incident == nil {
		http.Error(w, "incident not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incident)
}

// UpdateIncident handles PUT /incidents/{id}
func (h *IncidentHandler) UpdateIncident(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/incident/"):]
	if idStr == "" {
		http.Error(w, "missing incident id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid incident id", http.StatusBadRequest)
		return
	}

	var incident models.Incident
	if err := json.NewDecoder(r.Body).Decode(&incident); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	incident.ID = id // ensure the ID from the URL is used

	if err := h.Service.UpdateIncident(r.Context(), &incident); err != nil {
		log.Printf("failed to update incident: %v", err)
		http.Error(w, "failed to update incident", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incident)
}

// DeleteIncident handles DELETE /incidents/{id}
func (h *IncidentHandler) DeleteIncident(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/incident/"):]
	if idStr == "" {
		http.Error(w, "missing incident id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid incident id", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteIncident(r.Context(), id); err != nil {
		log.Printf("failed to delete incident: %v", err)
		http.Error(w, "failed to delete incident", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchIncidents handles GET /incidents?query=...
func (h *IncidentHandler) SearchIncidents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	limit := 10 // default limit
	if l := r.URL.Query().Get("limit"); l != "" {
		if li, err := strconv.Atoi(l); err == nil && li > 0 {
			limit = li
		}
	}

	incidents, err := h.Service.SearchIncidents(r.Context(), query, limit)
	if err != nil {
		log.Printf("failed to search incidents: %v", err)
		http.Error(w, "failed to search incidents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incidents)
}