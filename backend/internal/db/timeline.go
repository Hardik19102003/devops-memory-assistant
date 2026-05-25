package db

import (
	"fmt"
	"strings"
)

func BuildIncidentTimeline(chunks []RetrievedChunk) string {

	var builder strings.Builder

	if len(chunks) == 0 {
		return ""
	}

	builder.WriteString(
		fmt.Sprintf(
			"INCIDENT: %s\n\n",
			chunks[0].Error,
		),
	)

	builder.WriteString("TIMELINE:\n")

	for _, chunk := range chunks {

		builder.WriteString(
			fmt.Sprintf(
				"\nSTEP %d:\n%s\n",
				chunk.ChunkIndex+1,
				chunk.ChunkText,
			),
		)
	}

	return builder.String()
}