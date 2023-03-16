package gpt

import (
	"CallFrescoBot/pkg/messages"
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func GetResponse(question string) string {
	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		log.Fatalln(messages.MissingGptKey)
	}

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		log.Println(err)

		return messages.ErrorMsg
	}

	return resp.Choices[0].Message.Content
}
