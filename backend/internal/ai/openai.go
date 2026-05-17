package ai

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func GenerateSummary(issue string) (string, error) {

	client := openai.NewClient(
		os.Getenv("OPENAI_API_KEY"),
	)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,

			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `
You are a senior DevOps engineer.

Summarize incidents shortly and professionally.
`,
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: issue,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}