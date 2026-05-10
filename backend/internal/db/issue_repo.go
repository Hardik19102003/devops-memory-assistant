package db

import (
	"database/sql"
	"devops-memory-assistant/internal/models"
	"strings"

	"github.com/lib/pq"
)

func SaveIssue(issue models.Issue) error {

	_, err := DB.Exec(`
INSERT INTO issues (
	error,
	cause,
	fix,
	steps,
	tags
)
VALUES ($1, $2, $3, $4, $5)
`,
		issue.Error,
		strings.Join(issue.Causes, " | "),
		strings.Join(issue.Fixes, " | "),
		strings.Join(issue.DebugSteps, " | "),
		pq.Array(issue.Tags),
	)

	return err
}

func SearchIssue(query string) ([]models.Issue, error) {

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
			&issue.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Convert DB strings back into arrays
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

	// Convert DB string → arrays
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