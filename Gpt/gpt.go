package gpt

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	usageService "CallFrescoBot/pkg/service/usage"
	"CallFrescoBot/pkg/utils"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

var ImageSupportedModels = []string{
	openai.GPT4oMini,
	openai.GPT4o,
}

func GetResponse(bot *tg.BotAPI, update tg.Update, user *models.User, model string) ([]tg.Chattable, error) {
	if utils.GetEnvVar("GPT_API_KEY") == "" {
		return nil, errors.New(consts.ErrorMissingGptKey)
	}

	client := openai.NewClient(utils.GetEnvVar("GPT_API_KEY"))

	sendMsgErr := messageService.SendMsgToUser(update.Message.Chat.ID, utils.LocalizeSafe(consts.GptLoading))
	if sendMsgErr != nil {
		return nil, sendMsgErr
	}

	if update.Message.Photo != nil && isImageSupportedModel(model) {
		return handleImageMessage(bot, update, user, client, model)
	} else {
		return handleTextMessage(update, user, client, model)
	}
}

func isImageSupportedModel(model string) bool {
	for _, m := range ImageSupportedModels {
		if model == m {
			return true
		}
	}
	return false
}

func handleImageMessage(bot *tg.BotAPI, update tg.Update, user *models.User, client *openai.Client, model string) ([]tg.Chattable, error) {
	photo := update.Message.Photo
	fileID := photo[len(photo)-1].FileID
	file, err := bot.GetFile(tg.FileConfig{FileID: fileID})
	if err != nil {
		return nil, fmt.Errorf("error getting file: %w", err)
	}

	fileURL := file.Link(bot.Token)
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading image data: %w", err)
	}

	imgBase64 := base64.StdEncoding.EncodeToString(imageData)

	imageDataURL := fmt.Sprintf("data:image/jpeg;base64,%s", imgBase64)

	var multiContent []openai.ChatMessagePart

	if update.Message.Caption != "" {
		multiContent = append(multiContent, openai.ChatMessagePart{
			Type: openai.ChatMessagePartTypeText,
			Text: update.Message.Caption,
		})
	}

	multiContent = append(multiContent, openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeImageURL,
		ImageURL: &openai.ChatMessageImageURL{
			URL:    imageDataURL,
			Detail: openai.ImageURLDetailLow,
		},
	})

	request := openai.ChatCompletionRequest{
		Model:     model,
		MaxTokens: 1000,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:         openai.ChatMessageRoleUser,
				MultiContent: multiContent,
			},
		},
	}

	res, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return nil, fmt.Errorf("error getting response from GPT: %w", err)
	}

	return handleGptResponse(update, user, res)
}

func handleTextMessage(update tg.Update, user *models.User, client *openai.Client, model string) ([]tg.Chattable, error) {
	request := createRequest(user, update, model)

	res, err := getResponseFromGPT(client, request)
	if err != nil {
		return nil, fmt.Errorf("error getting response from GPT: %w", err)
	}

	return handleGptResponse(update, user, res)
}

func getResponseFromGPT(client *openai.Client, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	res, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return openai.ChatCompletionResponse{}, fmt.Errorf("error creating chat completion: %w", err)
	}
	return res, nil
}

func handleGptResponse(update tg.Update, user *models.User, res openai.ChatCompletionResponse) ([]tg.Chattable, error) {
	var responseContent string
	if len(res.Choices) > 0 {
		responseContent = res.Choices[0].Message.Content
	} else {
		return nil, errors.New("no response from GPT")
	}

	var inputMessage string
	if update.Message.Caption != "" {
		inputMessage = update.Message.Caption
	} else {
		inputMessage = update.Message.Text
	}

	err := messageService.CreateMessage(user.Id, inputMessage, responseContent, user.Mode)
	if err != nil {
		return nil, fmt.Errorf("error creating message: %w", err)
	}

	userMode := user.Mode
	if user.Dialog == 1 {
		userMode = userMode + 100
	}

	if update.Message.Voice != nil {
		userMode = userMode + 1000
	}

	err = usageService.SaveUsage(user, userMode)
	if err != nil {
		return nil, fmt.Errorf("error saving usage: %w", err)
	}

	text := responseContent
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
