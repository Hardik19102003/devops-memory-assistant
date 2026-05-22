package ai

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func GenerateSummary(issue string) (string, error) {

	key := os.Getenv("OPENAI_API_KEY")

	log.Println("🔑 OPENAI KEY PRESENT?", key != "")
	log.Println("🔑 KEY LENGTH:", len(key))

	client := openai.NewClient(key)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "You are a senior DevOps engineer. Summarize incidents shortly.",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: issue,
				},
			},
		},
	)

	if err != nil {
		log.Println("❌ OPENAI ERROR:", err)
		return "", err
	}

	log.Println("✅ OPENAI SUCCESS")
	return resp.Choices[0].Message.Content, nil
}