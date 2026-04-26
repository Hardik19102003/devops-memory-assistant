package db

import "devops-memory-assistant/internal/models"

func SaveIssue(issue models.Issue) error {
	query := "INSERT INTO issues (error, cause, fix) VALUES ($1, $2, $3)"
	_, err := DB.Exec(query, issue.Error, issue.Cause, issue.Fix)
	return err
}

func (d *Database) GetIssuesByError(errText string) ([]models.Issue, error) {
	rows, err := d.DB.Query(
		`SELECT id, error, cause, fix FROM issues WHERE error ILIKE '%' || $1 || '%'`,
		errText,
	)
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