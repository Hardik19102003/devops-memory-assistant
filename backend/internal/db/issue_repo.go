package db

import "devops-memory-assistant/internal/models"

func SaveIssue(issue models.Issue) error {
	query := "INSERT INTO issues (error, cause, fix) VALUES ($1, $2, $3)"
	_, err := DB.Exec(query, issue.Error, issue.Cause, issue.Fix)
	return err
}

func GetIssuesByError(search string) ([]models.Issue, error) {

	query := `
	SELECT id, error, cause, fix 
	FROM issues 
	WHERE error ILIKE '%' || $1 || '%'
	`

	rows, err := DB.Query(query, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []models.Issue

	for rows.Next() {
		var i models.Issue
		rows.Scan(&i.ID, &i.Error, &i.Cause, &i.Fix)
		issues = append(issues, i)
	}

	return issues, nil
}