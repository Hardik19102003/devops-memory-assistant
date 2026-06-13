package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"

	"devops-memory-assistant/internal/models"
)

// NullVector wraps pgvector.Vector to support NULL values.
type NullVector struct {
	Vector pgvector.Vector
	Valid  bool
}

func (nv *NullVector) Scan(value interface{}) error {
	if value == nil {
		nv.Vector = pgvector.Vector{}
		nv.Valid = false
		return nil
	}
	nv.Valid = true
	return nv.Vector.Scan(value)
}

func (nv NullVector) Value() (driver.Value, error) {
	if !nv.Valid {
		return nil, nil
	}
	return nv.Vector.Value()
}

const (
	tableIncidents = "incidents"
)

// SaveIncident saves a new incident to the database.
func SaveIncident(incident models.Incident) (int, error) {
	ctx := context.Background()
	query := fmt.Sprintf(`
		INSERT INTO %s (
			title, summary, symptoms, evidence, root_cause, resolution, prevention,
			commands_used, tags, severity, environment, services_affected, lessons_learned,
			raw_notes, embedding, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		) RETURNING id
	`, tableIncidents)

	var id int
	// Convert []float64 to []float32 for pgvector
	if incident.Embedding != nil {
		embedding32 := make([]float32, len(incident.Embedding))
		for i, v := range incident.Embedding {
			embedding32[i] = float32(v)
		}

		err := DB.QueryRowContext(ctx, query,
			incident.Title,
			incident.Summary,
			pq.Array(incident.Symptoms),
			pq.Array(incident.Evidence),
			pq.Array(incident.RootCause),
			pq.Array(incident.Resolution),
			pq.Array(incident.Prevention),
			pq.Array(incident.CommandsUsed),
			pq.Array(incident.Tags),
			incident.Severity,
			incident.Environment,
			pq.Array(incident.ServicesAffected),
			incident.LessonsLearned,
			incident.RawNotes,
			pgvector.NewVector(embedding32),
			time.Now(),
			time.Now(),
		).Scan(&id)

		if err != nil {
			return 0, fmt.Errorf("failed to save incident: %w", err)
		}

		return id, nil
	}

	// Handle nil embedding case
	err := DB.QueryRowContext(ctx, query,
		incident.Title,
		incident.Summary,
		pq.Array(incident.Symptoms),
		pq.Array(incident.Evidence),
		pq.Array(incident.RootCause),
		pq.Array(incident.Resolution),
		pq.Array(incident.Prevention),
		pq.Array(incident.CommandsUsed),
		pq.Array(incident.Tags),
		incident.Severity,
		incident.Environment,
		pq.Array(incident.ServicesAffected),
		incident.LessonsLearned,
		incident.RawNotes,
		nil,
		time.Now(),
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to save incident: %w", err)
	}

	return id, nil
}

// GetIncidentByID retrieves an incident by its ID.
func GetIncidentByID(id int) (*models.Incident, error) {
	ctx := context.Background()
	query := fmt.Sprintf(`
		SELECT
			id, title, summary, symptoms, evidence, root_cause, resolution, prevention,
			commands_used, tags, severity, environment, services_affected, lessons_learned,
			raw_notes, embedding, created_at, updated_at
		FROM %s
		WHERE id = $1
	`, tableIncidents)

	var incident models.Incident
	var createdAt, updatedAt time.Time
	var embedding NullVector
	err := DB.QueryRowContext(ctx, query, id).Scan(
		&incident.ID,
		&incident.Title,
		&incident.Summary,
		pq.Array(&incident.Symptoms),
		pq.Array(&incident.Evidence),
		pq.Array(&incident.RootCause),
		pq.Array(&incident.Resolution),
		pq.Array(&incident.Prevention),
		pq.Array(&incident.CommandsUsed),
		pq.Array(&incident.Tags),
		&incident.Severity,
		&incident.Environment,
		pq.Array(&incident.ServicesAffected),
		&incident.LessonsLearned,
		&incident.RawNotes,
		&embedding,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get incident: %w", err)
	}
	// Convert pgvector.Vector to []float64
	if embedding.Valid && len(embedding.Vector.Slice()) > 0 {
		slice := embedding.Vector.Slice()
		incident.Embedding = make([]float64, len(slice))
		for i, v := range slice {
			incident.Embedding[i] = float64(v)
		}
	}
	incident.CreatedAt = createdAt
	incident.UpdatedAt = updatedAt
	return &incident, nil
}

// UpdateIncident updates an existing incident.
func UpdateIncident(incident models.Incident) error {
	ctx := context.Background()
	query := fmt.Sprintf(`
		UPDATE %s SET
			title = $1,
			summary = $2,
			symptoms = $3,
			evidence = $4,
			root_cause = $5,
			resolution = $6,
			prevention = $7,
			commands_used = $8,
			tags = $9,
			severity = $10,
			environment = $11,
			services_affected = $12,
			lessons_learned = $13,
			raw_notes = $14,
			embedding = $15,
			updated_at = $16
		WHERE id = $17
	`, tableIncidents)

	_, err := DB.ExecContext(ctx, query,
		incident.Title,
		incident.Summary,
		pq.Array(incident.Symptoms),
		pq.Array(incident.Evidence),
		pq.Array(incident.RootCause),
		pq.Array(incident.Resolution),
		pq.Array(incident.Prevention),
		pq.Array(incident.CommandsUsed),
		pq.Array(incident.Tags),
		incident.Severity,
		incident.Environment,
		pq.Array(incident.ServicesAffected),
		incident.LessonsLearned,
		incident.RawNotes,
		func() interface{} {
			if incident.Embedding == nil {
				return nil
			}
			embedding32 := make([]float32, len(incident.Embedding))
			for i, v := range incident.Embedding {
				embedding32[i] = float32(v)
			}
			return pgvector.NewVector(embedding32)
		}(),
		time.Now(),
		incident.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update incident: %w", err)
	}
	return nil
}

// DeleteIncident deletes an incident by its ID.
func DeleteIncident(id int) error {
	ctx := context.Background()
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableIncidents)
	_, err := DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete incident: %w", err)
	}
	return nil
}

// SearchIncidents searches for incidents based on a query string (full-text or semantic).
// This is a placeholder; we'll implement semantic search using embeddings later.
func SearchIncidents(query string, limit int) ([]models.Incident, error) {
	ctx := context.Background()
	// For now, we'll do a simple text search on title and summary.
	// In the future, we'll use the embedding column for semantic search.
	q := fmt.Sprintf("SELECT id, title, summary, symptoms, evidence, root_cause, resolution, prevention, commands_used, tags, severity, environment, services_affected, lessons_learned, raw_notes, embedding, created_at, updated_at FROM %s WHERE title ILIKE $1 OR summary ILIKE $1 LIMIT $2", tableIncidents)
	rows, err := DB.QueryContext(ctx, q, "%"+query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search incidents: %w", err)
	}
	defer rows.Close()

	var incidents []models.Incident
	for rows.Next() {
		var incident models.Incident
		var createdAt, updatedAt time.Time
		var embedding NullVector
		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Summary,
			pq.Array(&incident.Symptoms),
			pq.Array(&incident.Evidence),
			pq.Array(&incident.RootCause),
			pq.Array(&incident.Resolution),
			pq.Array(&incident.Prevention),
			pq.Array(&incident.CommandsUsed),
			pq.Array(&incident.Tags),
			&incident.Severity,
			&incident.Environment,
			pq.Array(&incident.ServicesAffected),
			&incident.LessonsLearned,
			&incident.RawNotes,
			&embedding,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		// Convert pgvector.Vector to []float64
		if embedding.Valid && len(embedding.Vector.Slice()) > 0 {
			slice := embedding.Vector.Slice()
			incident.Embedding = make([]float64, len(slice))
			for i, v := range slice {
				incident.Embedding[i] = float64(v)
			}
		}
		incident.CreatedAt = createdAt
		incident.UpdatedAt = updatedAt
		incidents = append(incidents, incident)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over incident rows: %w", err)
	}
	return incidents, nil
}

// IncidentWithDistance wraps an incident with its pgvector cosine distance.
type IncidentWithDistance struct {
	Incident models.Incident
	Distance float64
}

// SearchIncidentsByEmbeddingWithDistance returns incidents ordered by similarity
// along with their cosine distance (0 = identical, 2 = opposite).
func SearchIncidentsByEmbeddingWithDistance(embedding []float32, limit int) ([]IncidentWithDistance, error) {
	ctx := context.Background()
	q := fmt.Sprintf(`
		SELECT id, title, summary, symptoms, evidence, root_cause, resolution, prevention,
		       commands_used, tags, severity, environment, services_affected, lessons_learned,
		       raw_notes, embedding, created_at, updated_at,
		       embedding <=> $1 AS distance
		FROM %s
		WHERE embedding IS NOT NULL
		ORDER BY distance ASC
		LIMIT $2
	`, tableIncidents)
	rows, err := DB.QueryContext(ctx, q, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search incidents by embedding: %w", err)
	}
	defer rows.Close()

	var results []IncidentWithDistance
	for rows.Next() {
		var incident models.Incident
		var createdAt, updatedAt time.Time
		var emb NullVector
		var distance float64
		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Summary,
			pq.Array(&incident.Symptoms),
			pq.Array(&incident.Evidence),
			pq.Array(&incident.RootCause),
			pq.Array(&incident.Resolution),
			pq.Array(&incident.Prevention),
			pq.Array(&incident.CommandsUsed),
			pq.Array(&incident.Tags),
			&incident.Severity,
			&incident.Environment,
			pq.Array(&incident.ServicesAffected),
			&incident.LessonsLearned,
			&incident.RawNotes,
			&emb,
			&createdAt,
			&updatedAt,
			&distance,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		if emb.Valid && len(emb.Vector.Slice()) > 0 {
			slice := emb.Vector.Slice()
			incident.Embedding = make([]float64, len(slice))
			for i, v := range slice {
				incident.Embedding[i] = float64(v)
			}
		}
		incident.CreatedAt = createdAt
		incident.UpdatedAt = updatedAt
		results = append(results, IncidentWithDistance{Incident: incident, Distance: distance})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over incident rows: %w", err)
	}
	return results, nil
}

// SearchIncidentsByEmbedding finds incidents semantically similar to the given embedding vector.
func SearchIncidentsByEmbedding(embedding []float32, limit int) ([]models.Incident, error) {
	ctx := context.Background()
	q := fmt.Sprintf(`
		SELECT id, title, summary, symptoms, evidence, root_cause, resolution, prevention,
		       commands_used, tags, severity, environment, services_affected, lessons_learned,
		       raw_notes, embedding, created_at, updated_at
		FROM %s
		WHERE embedding IS NOT NULL
		ORDER BY embedding <=> $1
		LIMIT $2
	`, tableIncidents)
	rows, err := DB.QueryContext(ctx, q, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search incidents by embedding: %w", err)
	}
	defer rows.Close()

	var incidents []models.Incident
	for rows.Next() {
		var incident models.Incident
		var createdAt, updatedAt time.Time
		var emb NullVector
		err := rows.Scan(
			&incident.ID,
			&incident.Title,
			&incident.Summary,
			pq.Array(&incident.Symptoms),
			pq.Array(&incident.Evidence),
			pq.Array(&incident.RootCause),
			pq.Array(&incident.Resolution),
			pq.Array(&incident.Prevention),
			pq.Array(&incident.CommandsUsed),
			pq.Array(&incident.Tags),
			&incident.Severity,
			&incident.Environment,
			pq.Array(&incident.ServicesAffected),
			&incident.LessonsLearned,
			&incident.RawNotes,
			&emb,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan incident: %w", err)
		}
		if emb.Valid && len(emb.Vector.Slice()) > 0 {
			slice := emb.Vector.Slice()
			incident.Embedding = make([]float64, len(slice))
			for i, v := range slice {
				incident.Embedding[i] = float64(v)
			}
		}
		incident.CreatedAt = createdAt
		incident.UpdatedAt = updatedAt
		incidents = append(incidents, incident)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over incident rows: %w", err)
	}
	return incidents, nil
}