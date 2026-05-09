package db

import "devops-memory-assistant/internal/models"
import "github.com/lib/pq"
import "database/sql"

func SaveIssue(issue models.Issue) error {
	query := "INSERT INTO issues (error, cause, fix, steps, tags) VALUES ($1, $2, $3, $4, $5)"
	_, err := DB.Exec(
	query,
	issue.Error,
	issue.Cause,
	issue.Fix,
	issue.Steps,
	pq.Array(issue.Tags), // ✅ important
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

		err := rows.Scan(
	&issue.ID,
	&issue.Error,
	&issue.Cause,
	&issue.Fix,
	&issue.Steps,
	pq.Array(&issue.Tags), // 🔥 important
	&issue.CreatedAt,
)
		if err != nil {
			return nil, err
		}

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
	err := row.Scan(
	&issue.ID,
	&issue.Error,
	&issue.Cause,
	&issue.Fix,
	&issue.Steps,
	pq.Array(&issue.Tags), // ✅ FIX
	&issue.CreatedAt,
)

	if err == sql.ErrNoRows {
	return nil, nil
}
if err != nil {
	return nil, err
}

	return &issue, nil
}

func DeleteIssue(id string) error {

	query := `DELETE FROM issues WHERE id = $1`

	_, err := DB.Exec(query, id)

	return err
}
