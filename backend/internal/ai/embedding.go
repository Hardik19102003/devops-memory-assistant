package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OllamaEmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func GenerateEmbedding(text string) ([]float32, error) {

	bodyMap := map[string]interface{}{
		"model":  "nomic-embed-text",
		"prompt": text,
	}

	bodyBytes, _ := json.Marshal(bodyMap)

	resp, err := http.Post(
		"http://localhost:11434/api/embeddings",
		"application/json",
		bytes.NewBuffer(bodyBytes),
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result OllamaEmbeddingResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Println("✅ Embedding generated:", len(result.Embedding))

	return result.Embedding, nil
}