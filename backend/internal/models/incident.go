package models

import (
	"time"
)

// Incident represents a structured troubleshooting entry
type Incident struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Summary     string    `json:"summary" db:"summary"`
	Symptoms    []string  `json:"symptoms" db:"symptoms"`
	Evidence    []string  `json:"evidence" db:"evidence"`
	RootCause   []string  `json:"root_cause" db:"root_cause"`
	Resolution  []string  `json:"resolution" db:"resolution"`
	Prevention  []string  `json:"prevention" db:"prevention"`
	CommandsUsed []string `json:"commands_used" db:"commands_used"`
	Tags        []string  `json:"tags" db:"tags"`
	Severity    string    `json:"severity" db:"severity"` // low, medium, high, critical
	Environment string    `json:"environment" db:"environment"`
	ServicesAffected []string `json:"services_affected" db:"services_affected"`
	LessonsLearned string  `json:"lessons_learned" db:"lessons_learned"`
	RawNotes    string    `json:"raw_notes" db:"raw_notes"`
	Embedding   []float64 `json:"embedding" db:"embedding"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// IncidentInput represents the input for creating/updating an incident
type IncidentInput struct {
	Title       string   `json:"title"`
	Summary     string   `json:"summary"`
	Symptoms    []string `json:"symptoms"`
	Evidence    []string `json:"evidence"`
	RootCause   []string `json:"root_cause"`
	Resolution  []string `json:"resolution"`
	Prevention  []string `json:"prevention"`
	CommandsUsed []string `json:"commands_used"`
	Tags        []string `json:"tags"`
	Severity    string   `json:"severity"`
	Environment string   `json:"environment"`
	ServicesAffected []string `json:"services_affected"`
	LessonsLearned string  `json:"lessons_learned"`
	RawNotes    string   `json:"raw_notes"`
}