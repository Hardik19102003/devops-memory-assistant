package service

import (
	"context"
	"fmt"
	"log"

	"devops-memory-assistant/internal/ai"
	"devops-memory-assistant/internal/db"
	"devops-memory-assistant/internal/models"
)

// IncidentService handles incident related business logic.
type IncidentService struct{}

// NewIncidentService creates a new IncidentService.
func NewIncidentService() *IncidentService {
	return &IncidentService{}
}

// ExtractIncidentFromNotes extracts structured incident data from raw notes.
// If similar past incidents exist in the DB, they are injected into the LLM prompt (RAG).
func (s *IncidentService) ExtractIncidentFromNotes(ctx context.Context, rawNotes string) (*models.IncidentInput, error) {
	// 1. Generate embedding for the raw notes
	embedding, err := ai.GenerateEmbedding(rawNotes)
	if err != nil {
		log.Printf("[RAG] embedding generation failed: %v", err)
	}

	// 2. Find semantically similar past incidents
	var related []models.Incident
	if err == nil && embedding != nil {
		related, err = db.SearchIncidentsByEmbedding(embedding, 3)
		if err != nil {
			log.Printf("[RAG] semantic search failed: %v", err)
		}
		if len(related) > 0 {
			log.Printf("[RAG] found %d related incidents", len(related))
		}
	}

	// 3. Extract with context (RAG-enhanced prompt)
	incidentInput, err := ai.ExtractIncidentWithContext(rawNotes, related)
	if err != nil {
		return nil, err
	}
	incidentInput.RawNotes = rawNotes
	return incidentInput, nil
}

// SaveIncident saves an incident to the database.
func (s *IncidentService) SaveIncident(ctx context.Context, incidentInput *models.IncidentInput) (*models.Incident, error) {
	// Convert IncidentInput to Incident model.
	incident := models.Incident{
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
		}

	// Generate embedding from raw notes and store it so future RAG searches work.
	embedding, err := ai.GenerateEmbedding(incidentInput.RawNotes)
	if err == nil && embedding != nil {
		incident.Embedding = make([]float64, len(embedding))
		for i, v := range embedding {
			incident.Embedding[i] = float64(v)
		}
	} else if err != nil {
		log.Printf("[SAVE] embedding generation failed: %v", err)
	}

	// Save the incident using the db package.
	id, err := db.SaveIncident(incident)
	if err != nil {
		return nil, err
	}
	incident.ID = id

	// Fetch the saved incident to get the generated fields (like timestamps, embedding).
	savedIncident, err := db.GetIncidentByID(id)
	if err != nil {
		return nil, err
	}

	return savedIncident, nil
}

// GetIncidentByID retrieves an incident by its ID.
func (s *IncidentService) GetIncidentByID(ctx context.Context, id int) (*models.Incident, error) {
	return db.GetIncidentByID(id)
}

// UpdateIncident updates an existing incident.
func (s *IncidentService) UpdateIncident(ctx context.Context, incident *models.Incident) error {
	return db.UpdateIncident(*incident)
}

// DeleteIncident deletes an incident by its ID.
func (s *IncidentService) DeleteIncident(ctx context.Context, id int) error {
	return db.DeleteIncident(id)
}

// FindSimilarIncidents performs semantic search using embeddings and returns
// the top-N most similar past incidents with similarity percentages.
func (s *IncidentService) FindSimilarIncidents(ctx context.Context, rawNotes string, limit int) ([]models.SimilarIncident, error) {
	embedding, err := ai.GenerateEmbedding(rawNotes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	results, err := db.SearchIncidentsByEmbeddingWithDistance(embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search incidents: %w", err)
	}

	var similar []models.SimilarIncident
	for _, r := range results {
		score := 1.0 - r.Distance // cosine distance → similarity
		if score < 0 {
			score = 0
		}
		similar = append(similar, models.SimilarIncident{
			ID:         r.Incident.ID,
			Title:      r.Incident.Title,
			Summary:    r.Incident.Summary,
			Similarity: score,
		})
	}
	return similar, nil
}

// SearchIncidents searches for incidents based on a query.
func (s *IncidentService) SearchIncidents(ctx context.Context, query string, limit int) ([]models.Incident, error) {
	return db.SearchIncidents(query, limit)
}