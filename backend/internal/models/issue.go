package models

import (
	"time"
)

type Issue struct {
	ID        int            `json:"id"`
	Error     string         `json:"error"`
	Cause     string         `json:"cause"`
	Fix       string         `json:"fix"`
	Steps     string `json:"steps"`
	Tags      []string       `json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
}
