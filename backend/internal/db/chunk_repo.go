package db

import (
	"strings"

	"devops-memory-assistant/internal/ai"

	"github.com/pgvector/pgvector-go"
)

func SaveChunks(issueID int, document string) error {

	// Split by double newlines (paragraphs)
	rawChunks := strings.Split(document, "\n\n")

	var chunks []string

	for _, chunk := range rawChunks {

		cleaned := strings.TrimSpace(chunk)

		if cleaned == "" {
			continue
		}

		chunks = append(chunks, cleaned)
	}

	for index, chunk := range chunks {

		embedding, err := ai.GenerateEmbedding(chunk)

		if err != nil {
			continue
		}

		_, err = DB.Exec(`
INSERT INTO issue_chunks (
	issue_id,
	chunk_index,
	content,
	embedding
)
VALUES ($1, $2, $3, $4)
`,
			issueID,
			index,
			chunk,
			pgvector.NewVector(embedding),
		)

		if err != nil {
			return err
		}
	}

	return nil
}