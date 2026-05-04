package models

import (
	"database/sql"
	"time"
)

type Issue struct {
	ID        int            `json:"id"`
	Error     string         `json:"error"`
	Cause     string         `json:"cause"`
	Fix       string         `json:"fix"`
	Steps     sql.NullString `json:"steps"`
	Tags      []string       `json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
}
