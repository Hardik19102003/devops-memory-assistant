package models

type Issue struct {
	ID int `json:"id"`

	// Main issue title
	Error string `json:"error"`

	// Multiple possible causes
	Causes []string `json:"causes"`

	// Multiple possible fixes
	Fixes []string `json:"fixes"`

	// Step-by-step debugging flow
	DebugSteps []string `json:"debug_steps"`

	// Helpful commands
	Commands []string `json:"commands"`

	// Categorization
	Tags []string `json:"tags"`

	// Future AI recommendations
	RelatedIssues []string `json:"related_issues"`

	// Docs / URLs
	References []string `json:"references"`

	CreatedAt string `json:"created_at"`
}