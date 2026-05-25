package ai

import "strings"

func ChunkDocument(document string, chunkSize int) []string {

	var chunks []string

	// Clean document
	document = strings.TrimSpace(document)

	// Split by paragraphs first
	paragraphs := strings.Split(document, "\n")

	var currentChunk string

	for _, para := range paragraphs {

		para = strings.TrimSpace(para)

		if para == "" {
			continue
		}

		// If adding paragraph exceeds chunk size
		if len(currentChunk)+len(para) > chunkSize {

			chunks = append(chunks, currentChunk)

			currentChunk = para

		} else {

			if currentChunk == "" {

				currentChunk = para

			} else {

				currentChunk += "\n" + para
			}
		}
	}

	// Add last chunk
	if currentChunk != "" {

		chunks = append(chunks, currentChunk)
	}

	return chunks
}