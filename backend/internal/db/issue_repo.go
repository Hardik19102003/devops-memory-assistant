package db

import "devops-memory-assistant/internal/models"

func SaveIssue(issue models.Issue) error {
	query := "INSERT INTO issues (error, cause, fix) VALUES ($1, $2, $3)"
	_, err := DB.Exec(query, issue.Error, issue.Cause, issue.Fix)
	return err
}

func SearchIssue(query string) ([]models.Issue, error) {
	issues := []models.Issue{}

	searchTerm := "%" + query + "%"

	rows, err := DB.Query(`
		SELECT error, cause, fix, created_at
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
			&issue.Error,
			&issue.Cause,
			&issue.Fix,
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
		SELECT error, cause, fix, created_at
		FROM issues
		WHERE error ILIKE $1
		ORDER BY created_at DESC
		LIMIT 1
	`, searchTerm)

	var issue models.Issue
	err := row.Scan(
		&issue.Error,
		&issue.Cause,
		&issue.Fix,
		&issue.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &issue, nil
}