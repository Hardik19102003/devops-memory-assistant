package db

import (
	"database/sql"
	"strings"

	"devops-memory-assistant/internal/ai"
	"devops-memory-assistant/internal/models"

	"github.com/lib/pq"
	"github.com/pgvector/pgvector-go"
)

func SaveIssue(issue models.Issue) error {

	var vector interface{}

	// Better embedding source
	contentForEmbedding := `
	Issue: ` + issue.Error + `

	Document:
	` + issue.Document

	embedding, err := ai.GenerateEmbedding(contentForEmbedding)

	if err == nil {
		vector = pgvector.NewVector(embedding)
	}

	_, err = DB.Exec(`
INSERT INTO issues (
	error,
	cause,
	fix,
	steps,
	tags,
	document,
	embedding
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`,
		issue.Error,
		strings.Join(issue.Causes, " | "),
		strings.Join(issue.Fixes, " | "),
		strings.Join(issue.DebugSteps, " | "),
		pq.Array(issue.Tags),
		issue.Document,
		vector,
	)

	return err
}

func SearchIssue(query string) ([]models.Issue, error) {

	// Try semantic search first
	embedding, err := ai.GenerateEmbedding(query)

	if err != nil {

		// Fallback to old search
		return SearchIssueFallback(query)
	}

	issues := []models.Issue{}

	rows, err := DB.Query(`
	SELECT
		id,
		error,
		cause,
		fix,
		COALESCE(steps, ''),
		COALESCE(tags, '{}'::text[]),
		COALESCE(document, ''),
		created_at
	FROM issues
	ORDER BY embedding <-> $1
	LIMIT 5;
	`,
		pgvector.NewVector(embedding),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var issue models.Issue

		var causes string
		var fixes string
		var steps string

		err := rows.Scan(
			&issue.ID,
			&issue.Error,
			&causes,
			&fixes,
			&steps,
			pq.Array(&issue.Tags),
			&issue.Document,
			&issue.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		issue.Causes = strings.Split(causes, " | ")
		issue.Fixes = strings.Split(fixes, " | ")
		issue.DebugSteps = strings.Split(steps, " | ")

		issues = append(issues, issue)
	}

	return issues, nil
}

func SearchIssueFallback(query string) ([]models.Issue, error) {

	issues := []models.Issue{}

	searchTerm := "%" + query + "%"

	rows, err := DB.Query(`
	SELECT
		id,
		error,
		cause,
		fix,
		COALESCE(steps, ''),
		COALESCE(tags, '{}'::text[]),
		COALESCE(document, ''),
		created_at
	FROM issues
	WHERE error ILIKE $1
	   OR cause ILIKE $1
	   OR fix ILIKE $1
	ORDER BY created_at DESC
	LIMIT 5;
	`, searchTerm)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		var issue models.Issue

		var causes string
		var fixes string
		var steps string

		err := rows.Scan(
			&issue.ID,
			&issue.Error,
			&causes,
			&fixes,
			&steps,
			pq.Array(&issue.Tags),
			&issue.Document,
			&issue.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		issue.Causes = strings.Split(causes, " | ")
		issue.Fixes = strings.Split(fixes, " | ")
		issue.DebugSteps = strings.Split(steps, " | ")

		issues = append(issues, issue)
	}

	return issues, nil
}

func FindSimilarIssue(query string) (*models.Issue, error) {

	searchTerm := "%" + query + "%"

	row := DB.QueryRow(`
	SELECT 
		id,
		error,
		cause,
		fix,
		COALESCE(steps, ''),
		COALESCE(tags, '{}'::text[]),
		COALESCE(document, ''),
		created_at
	FROM issues
	WHERE error ILIKE $1
	ORDER BY created_at DESC
	LIMIT 1
`, searchTerm)

	var issue models.Issue

	var causes string
	var fixes string
	var steps string

	err := row.Scan(
		&issue.ID,
		&issue.Error,
		&causes,
		&fixes,
		&steps,
		pq.Array(&issue.Tags),
		&issue.Document,
		&issue.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	issue.Causes = strings.Split(causes, " | ")
	issue.Fixes = strings.Split(fixes, " | ")
	issue.DebugSteps = strings.Split(steps, " | ")

	return &issue, nil
}

func DeleteIssue(id string) error {

	query := `DELETE FROM issues WHERE id = $1`

	_, err := DB.Exec(query, id)

	return err
}