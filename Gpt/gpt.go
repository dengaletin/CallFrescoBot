package gpt

import (
	"CallFrescoBot/pkg/messages"
	"context"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
)

func GetResponse(question string) string {
	apiKey := os.Getenv("GPT_API_KEY")
	if apiKey == "" {
		log.Fatalln(messages.MissingGptKey)
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	var response string

	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt:      []string{question},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(res *gpt3.CompletionResponse) {
		response += res.Choices[0].Text
	})

	if err != nil {
		log.Println(err)

		return messages.ErrorMsg
	}
	return response
}
