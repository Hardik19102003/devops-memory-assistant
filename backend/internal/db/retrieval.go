package db

import (
	"strings"

	"devops-memory-assistant/internal/ai"

	"github.com/pgvector/pgvector-go"
)

type RelatedIncident struct {
	IssueID  int
	Error    string
	Distance float64
}

type RetrievedChunk struct {
	IssueID    int
	Error      string
	ChunkText  string
	ChunkIndex int
	Distance   float64
}

func FindRelatedIncidents(query string) ([]RelatedIncident, error) {

	embedding, err := ai.GenerateEmbedding(query)

	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
SELECT
	issues.id,
	issues.error,
	MIN(issue_chunks.embedding <-> $1) AS distance
FROM issue_chunks
JOIN issues
	ON issues.id = issue_chunks.issue_id
GROUP BY issues.id, issues.error
ORDER BY distance ASC
LIMIT 3;
`,
		pgvector.NewVector(embedding),
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []RelatedIncident

	for rows.Next() {

		var incident RelatedIncident

		err := rows.Scan(
			&incident.IssueID,
			&incident.Error,
			&incident.Distance,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, incident)
	}

	return results, nil
}

func GetIncidentChunks(issueID int) ([]RetrievedChunk, error) {

	rows, err := DB.Query(`
SELECT
	issues.id,
	issues.error,
	issue_chunks.content,
	issue_chunks.chunk_index
FROM issue_chunks
JOIN issues
	ON issues.id = issue_chunks.issue_id
WHERE issues.id = $1
ORDER BY issue_chunks.chunk_index ASC;
`,
		issueID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var chunks []RetrievedChunk

	for rows.Next() {

		var chunk RetrievedChunk

		err := rows.Scan(
			&chunk.IssueID,
			&chunk.Error,
			&chunk.ChunkText,
			&chunk.ChunkIndex,
		)

		if err != nil {
			return nil, err
		}

		chunk.ChunkText = strings.TrimSpace(chunk.ChunkText)

		chunks = append(chunks, chunk)
	}

	return chunks, nil
}