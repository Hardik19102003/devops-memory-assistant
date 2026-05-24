package models

type Issue struct {
	ID int `json:"id"`

	// Incident title
	Error string `json:"error"`

	// Structured fields
	Causes []string `json:"causes"`
	Fixes []string `json:"fixes"`
	DebugSteps []string `json:"debug_steps"`

	// Extra operational metadata
	Commands []string `json:"commands"`
	Tags []string `json:"tags"`

	// Long-form markdown incident doc
	Document string `json:"document"`

	// Future AI features
	RelatedIssues []string `json:"related_issues"`
	References []string `json:"references"`

	CreatedAt string `json:"created_at"`
}