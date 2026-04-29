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
		SELECT error, cause, fix
		FROM issues
		WHERE error ILIKE $1
		   OR cause ILIKE $1
		   OR fix ILIKE $1
	`, searchTerm)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var issue models.Issue
		err := rows.Scan(&issue.Error, &issue.Cause, &issue.Fix)
		if err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
