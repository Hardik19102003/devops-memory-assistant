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

	// 🔥 Create text for embedding
	text := issue.Error + " " +
		strings.Join(issue.Causes, " ") + " " +
		strings.Join(issue.Fixes, " ")

	// 🤖 Generate embedding using Ollama
	embedding, err := ai.GenerateEmbedding(text)

	if err != nil {
		return err
	}

	// 💾 Save into PostgreSQL
	_, err = DB.Exec(`
INSERT INTO issues (
	error,
	cause,
	fix,
	steps,
	tags,
	embedding
)
VALUES ($1, $2, $3, $4, $5, $6)
`,
		issue.Error,
		strings.Join(issue.Causes, " | "),
		strings.Join(issue.Fixes, " | "),
		strings.Join(issue.DebugSteps, " | "),
		pq.Array(issue.Tags),
		pgvector.NewVector(embedding),
	)

	return err
}

func SearchIssue(query string) ([]models.Issue, error) {

	issues := []models.Issue{}

	// Generate embedding for search query
	embedding, err := ai.GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
	SELECT
		id,
		error,
		cause,
		fix,
		COALESCE(steps, ''),
		COALESCE(tags, '{}'::text[]),
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