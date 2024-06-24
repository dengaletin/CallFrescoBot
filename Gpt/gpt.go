package gpt

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/service/message"
	usageService "CallFrescoBot/pkg/service/usage"
	"CallFrescoBot/pkg/utils"
	"context"
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

func getResponseFromGPT(client *openai.Client, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	res, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("error creating chat completion: %w", err)
	}
	return res, nil
}

func handleGptResponse(update tg.Update, user *models.User, res openai.ChatCompletionResponse) ([]tg.Chattable, error) {
	err := messageService.CreateMessage(user.Id, update.Message.Text, res.Choices[0].Message.Content, user.Mode)
	if err != nil {
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	userMode := user.Mode
	if user.Dialog == 1 {
		userMode = userMode + 100
	}

	err = usageService.SaveUsage(user, userMode)
	if err != nil {
		return nil, fmt.Errorf("error saving usage: %w", err)
	}

	text := res.Choices[0].Message.Content
	parts := splitMessage(text, 4095)

	var messages []tg.Chattable
	for _, part := range parts {
		message := tg.NewMessage(update.Message.Chat.ID, part)
		message.ReplyToMessageID = update.Message.MessageID
		messages = append(messages, message)
	}

	return messages, nil
}

func splitMessage(text string, maxLength int) []string {
	var parts []string
	for len(text) > maxLength {
		splitIndex := maxLength
		if len(text) > maxLength {
			splitIndex = strings.LastIndex(text[:maxLength], " ")
			if splitIndex == -1 {
				splitIndex = maxLength
			}
		}
		parts = append(parts, text[:splitIndex])
		text = text[splitIndex:]
	}
	parts = append(parts, text)

	return parts
}

func GetResponse(update tg.Update, user *models.User, model string) ([]tg.Chattable, error) {
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
		contextMessagesLimit := consts.ContextMessagesLimit

		if user.Id == consts.AdminUserId {
			contextMessagesLimit = consts.AdminContextMessagesLimit
		}

		userMessages, err := messageService.GetMessagesByUser(user, contextMessagesLimit, user.Mode)
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
