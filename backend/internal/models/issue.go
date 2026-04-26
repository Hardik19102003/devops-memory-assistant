package models

type Issue struct {
	ID    int    `json:"id"`
	Error string `json:"error"`
	Cause string `json:"cause"`
	Fix   string `json:"fix"`
}