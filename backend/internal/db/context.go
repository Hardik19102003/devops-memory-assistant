package db

import "strings"

func BuildMultiIncidentContext(issueIDs []int) (string, error) {

	var builder strings.Builder

	for _, issueID := range issueIDs {

		chunks, err := GetIncidentChunks(issueID)

		if err != nil {
			continue
		}

		if len(chunks) == 0 {
			continue
		}

		builder.WriteString("\n=================================\n")
		builder.WriteString("RELATED INCIDENT\n")
		builder.WriteString("=================================\n\n")

		builder.WriteString(
			BuildIncidentTimeline(chunks),
		)

		builder.WriteString("\n\n")
	}

	return builder.String(), nil
}