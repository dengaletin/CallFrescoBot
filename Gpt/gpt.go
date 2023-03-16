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

	// ctx := context.Background()
	// client := gpt3.NewClient(apiKey)

	// var response string

	// err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
	// 	Prompt:      []string{question},
	// 	MaxTokens:   gpt3.IntPtr(3000),
	// 	Temperature: gpt3.Float32Ptr(0),
	// }, func(res *gpt3.CompletionResponse) {
	// 	response += res.Choices[0].Text
	// })

	// if err != nil {
	// 	log.Println(err)

	// 	return messages.ErrorMsg
	// }
	// return response

	client := openai.NewClient("your token")
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
