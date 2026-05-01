package models

import "time"

type Issue struct {
	Error     string    `json:"error"`
	Cause     string    `json:"cause"`
	Fix       string    `json:"fix"`
	CreatedAt time.Time `json:"created_at"`
}