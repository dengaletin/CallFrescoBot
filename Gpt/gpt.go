package gpt

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/service/message"
	"CallFrescoBot/pkg/utils"
	"context"
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"log"
)

func getResponseFromGPT(client *openai.Client, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	res, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("error creating chat completion: %w", err)
	}
	return res, nil
}

func handleGptResponse(update tg.Update, user *models.User, res openai.ChatCompletionResponse) (tg.Chattable, error) {
	err := messageService.CreateMessage(user.Id, update.Message.Text, res.Choices[0].Message.Content, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	text := res.Choices[0].Message.Content

	message := tg.NewMessage(update.Message.Chat.ID, text)
	message.ReplyToMessageID = update.Message.MessageID

	return message, nil
}

func GetResponse(update tg.Update, user *models.User, model string) (tg.Chattable, error) {
	if utils.GetEnvVar("GPT_API_KEY") == "" {
		return nil, errors.New(consts.ErrorMissingGptKey)
	}

	client := openai.NewClient(utils.GetEnvVar("GPT_API_KEY"))

	sendMsgErr := messageService.SendMsgToUser(update.Message.Chat.ID, utils.LocalizeSafe(consts.GptLoading))
	if sendMsgErr != nil {
		return nil, sendMsgErr
	}

	request := createRequest(user, update, model)

	resp, err := getResponseFromGPT(client, request)
	if err != nil {
		return nil, fmt.Errorf("error getting response from GPT: %w", err)
	}

	return handleGptResponse(update, user, resp)
}

func createRequest(user *models.User, update tg.Update, model string) openai.ChatCompletionRequest {
	var messages []openai.ChatCompletionMessage

	if user.Dialog == 1 {
		userMessages, err := messageService.GetMessagesByUser(user, 15, 0)
		if err != nil {
			log.Printf("error getting messages by user: %v", err)
		}

		for _, userMessage := range userMessages {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: userMessage.Message,
			}, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: userMessage.Response,
			})
		}
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: update.Message.Text,
	})

	return openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}
}
