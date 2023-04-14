package gpt

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	"CallFrescoBot/pkg/utils"
	"context"
	openai "github.com/sashabaranov/go-openai"
	"log"
)

var apiKey string

func init() {
	apiKey = utils.GetEnvVar("GPT_API_KEY")
}

func GetResponse(question string, user *models.User) string {
	if apiKey == "" {
		log.Fatalln(consts.MissingGptKey)
	}

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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

		return consts.ErrorMsg
	}

	err = messageService.CreateMessage(user.Id, question, resp.Choices[0].Message.Content)
	if err != nil {
		return consts.ErrorMsg
	}

	return resp.Choices[0].Message.Content
}
