package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"

	"devops-memory-assistant/internal/models"
)

func TestSaveIncident(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB variable for testing
	originalDB := DB
	DB = db
	defer func() { DB = originalDB }()

	// Mock the query execution
	mock.ExpectQuery("INSERT INTO incidents").
		WithArgs(
			"Test Title",
			"Test Summary",
			pq.Array([]string{"symptom1"}),
			pq.Array([]string{"evidence1"}),
			pq.Array([]string{"cause1"}),
			pq.Array([]string{"fix1"}),
			pq.Array([]string{"prevention1"}),
			pq.Array([]string{"cmd1"}),
			pq.Array([]string{"tag1"}),
			"medium",
			"test-env",
			pq.Array([]string{"service1"}),
			"lessons learned",
			"raw notes",
			pgvector.NewVector([]float32{0.1, 0.2, 0.3}),
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Test incident
	incident := models.Incident{
		Title:          "Test Title",
		Summary:        "Test Summary",
		Symptoms:       []string{"symptom1"},
		Evidence:       []string{"evidence1"},
		RootCause:      []string{"cause1"},
		Resolution:     []string{"fix1"},
		Prevention:     []string{"prevention1"},
		CommandsUsed:   []string{"cmd1"},
		Tags:           []string{"tag1"},
		Severity:       "medium",
		Environment:    "test-env",
		ServicesAffected: []string{"service1"},
		LessonsLearned: "lessons learned",
		RawNotes:       "raw notes",
		Embedding:      []float64{0.1, 0.2, 0.3},
	}

	// Call function
	id, err := SaveIncident(incident)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetIncidentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB variable for testing
	originalDB := DB
	DB = db
	defer func() { DB = originalDB }()

	// Mock the query execution
	rows := sqlmock.NewRows([]string{
		"id", "title", "summary", "symptoms", "evidence", "root_cause", "resolution", "prevention",
		"commands_used", "tags", "severity", "environment", "services_affected", "lessons_learned",
		"raw_notes", "embedding", "created_at", "updated_at",
	}).
		AddRow(
			1,
			"Test Title",
			"Test Summary",
			pq.Array([]string{"symptom1"}),
			pq.Array([]string{"evidence1"}),
			pq.Array([]string{"cause1"}),
			pq.Array([]string{"fix1"}),
			pq.Array([]string{"prevention1"}),
			pq.Array([]string{"cmd1"}),
			pq.Array([]string{"tag1"}),
			"medium",
			"test-env",
			pq.Array([]string{"service1"}),
			"lessons learned",
			"raw notes",
			pgvector.NewVector([]float32{0.1, 0.2, 0.3}),
			time.Now(),
			time.Now(),
		)

	mock.ExpectQuery("SELECT.*FROM incidents WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	// Call function
	incident, err := GetIncidentByID(1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, incident)
	assert.Equal(t, 1, incident.ID)
	assert.Equal(t, "Test Title", incident.Title)
	assert.Equal(t, []string{"symptom1"}, incident.Symptoms)
	assert.Equal(t, []string{"evidence1"}, incident.Evidence)
	assert.Equal(t, []string{"cause1"}, incident.RootCause)
	assert.Equal(t, []string{"fix1"}, incident.Resolution)
	assert.Equal(t, []string{"prevention1"}, incident.Prevention)
	assert.Equal(t, []string{"cmd1"}, incident.CommandsUsed)
	assert.Equal(t, []string{"tag1"}, incident.Tags)
	assert.Equal(t, "medium", incident.Severity)
	assert.Equal(t, "test-env", incident.Environment)
	assert.Equal(t, []string{"service1"}, incident.ServicesAffected)
	assert.Equal(t, "lessons learned", incident.LessonsLearned)
	assert.Equal(t, "raw notes", incident.RawNotes)
	expectedEmbedding := []float64{float64(float32(0.1)), float64(float32(0.2)), float64(float32(0.3))}
	assert.Equal(t, expectedEmbedding, incident.Embedding)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetIncidentByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB variable for testing
	originalDB := DB
	DB = db
	defer func() { DB = originalDB }()

	// Mock returning no rows
	mock.ExpectQuery("SELECT.*FROM incidents WHERE id = \\$1").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	// Call function
	incident, err := GetIncidentByID(999)

	// Assertions
	assert.NoError(t, err)
	assert.Nil(t, incident)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateIncident(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB variable for testing
	originalDB := DB
	DB = db
	defer func() { DB = originalDB }()

	// Mock the update execution
	mock.ExpectExec("UPDATE incidents SET").
		WithArgs(
			"Updated Title",
			"Updated Summary",
			pq.Array([]string{"updated-symptom"}),
			pq.Array([]string{"updated-evidence"}),
			pq.Array([]string{"updated-cause"}),
			pq.Array([]string{"updated-fix"}),
			pq.Array([]string{"updated-prevention"}),
			pq.Array([]string{"updated-cmd"}),
			pq.Array([]string{"updated-tag"}),
			"high",
			"updated-env",
			pq.Array([]string{"updated-service"}),
			"updated lessons",
			"updated raw notes",
			pgvector.NewVector([]float32{0.4, 0.5, 0.6}),
			sqlmock.AnyArg(),
			1,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Test incident
	incident := models.Incident{
		ID:             1,
		Title:          "Updated Title",
		Summary:        "Updated Summary",
		Symptoms:       []string{"updated-symptom"},
		Evidence:       []string{"updated-evidence"},
		RootCause:      []string{"updated-cause"},
		Resolution:     []string{"updated-fix"},
		Prevention:     []string{"updated-prevention"},
		CommandsUsed:   []string{"updated-cmd"},
		Tags:           []string{"updated-tag"},
		Severity:       "high",
		Environment:    "updated-env",
		ServicesAffected: []string{"updated-service"},
		LessonsLearned: "updated lessons",
		RawNotes:       "updated raw notes",
		Embedding:      []float64{0.4, 0.5, 0.6},
	}

	// Call function
	err = UpdateIncident(incident)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteIncident(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB variable for testing
	originalDB := DB
	DB = db
	defer func() { DB = originalDB }()

	// Mock the delete execution
	mock.ExpectExec("DELETE FROM incidents WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call function
	err = DeleteIncident(1)

	// Assertions
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchIncidents(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Set the global DB variable for testing
	originalDB := DB
	DB = db
	defer func() { DB = originalDB }()

	// Mock the query execution
	rows := sqlmock.NewRows([]string{
		"id", "title", "summary", "symptoms", "evidence", "root_cause", "resolution", "prevention",
		"commands_used", "tags", "severity", "environment", "services_affected", "lessons_learned",
		"raw_notes", "embedding", "created_at", "updated_at",
	}).
		AddRow(
			1,
			"Test Incident",
			"Test Summary",
			pq.Array([]string{"symptom1"}),
			pq.Array([]string{"evidence1"}),
			pq.Array([]string{"cause1"}),
			pq.Array([]string{"fix1"}),
			pq.Array([]string{"prevention1"}),
			pq.Array([]string{"cmd1"}),
			pq.Array([]string{"tag1"}),
			"medium",
			"test-env",
			pq.Array([]string{"service1"}),
			"lessons learned",
			"raw notes",
			pgvector.NewVector([]float32{0.1, 0.2, 0.3}),
			time.Now(),
			time.Now(),
		)

	mock.ExpectQuery("SELECT.*FROM incidents WHERE title ILIKE \\$1 OR summary ILIKE \\$1 LIMIT \\$2").
		WithArgs("%test%", 10).
		WillReturnRows(rows)

	// Call function
	incidents, err := SearchIncidents("test", 10)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, incidents, 1)
	assert.Equal(t, 1, incidents[0].ID)
	assert.Equal(t, "Test Incident", incidents[0].Title)
	assert.Equal(t, []string{"symptom1"}, incidents[0].Symptoms)
	assert.NoError(t, mock.ExpectationsWereMet())
}