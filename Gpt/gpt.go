package gpt

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	"CallFrescoBot/pkg/utils"
	"context"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

var apiKey string

func init() {
	apiKey = utils.GetEnvVar("GPT_API_KEY")
}

func GetResponse(update tg.Update, user *models.User) (tg.Chattable, error) {
	if apiKey == "" {
		return nil, errors.New(consts.MissingGptKey)
	}

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: update.Message.Text,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	err = messageService.CreateMessage(user.Id, update.Message.Text, resp.Choices[0].Message.Content)
	if err != nil {
		return nil, err
	}

	message := tg.NewMessage(update.Message.Chat.ID, resp.Choices[0].Message.Content)
	message.ReplyToMessageID = update.Message.MessageID

	return message, nil
}
