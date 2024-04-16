package Dalle3

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	usageService "CallFrescoBot/pkg/service/usage"
	"CallFrescoBot/pkg/utils"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

var apiKey = utils.GetEnvVar("GPT_API_KEY")

func GetResponse(update tg.Update, user *models.User) (tg.Chattable, error) {
	if apiKey == "" {
		return nil, errors.New(consts.ErrorMissingGptKey)
	}

	sendMsgErr := messageService.SendMsgToUser(update.Message.Chat.ID, utils.LocalizeSafe(consts.DalleLoading))
	if sendMsgErr != nil {
		return nil, sendMsgErr
	}

	return getImageResponse(update, user)
}

func getImageResponse(update tg.Update, user *models.User) (tg.Chattable, error) {
	c := openai.NewClient(apiKey)
	ctx := context.Background()
	imgBytes, err := getImageBytes(ctx, c, update, user)

	if err != nil {
		return nil, err
	}

	return getMessage(update, imgBytes, user)
}

func getImageBytes(
	ctx context.Context,
	openaiClient *openai.Client,
	update tg.Update,
	user *models.User,
) ([]byte, error) {
	reqBase64 := prepareRequest(update.Message.Text, user)
	respBase64, err := openaiClient.CreateImage(ctx, reqBase64)

	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return nil, err
	}

	return base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
}

func prepareRequest(promptText string, user *models.User) openai.ImageRequest {
	var messages []string

	if user.Dialog == 1 {
		userMessages, err := messageService.GetMessagesByUser(user, consts.ContextMessagesLimit, 1)
		if err != nil {
			log.Printf("error getting messages by user: %v", err)
		} else {
			for _, userMessage := range userMessages {
				messages = append(messages, userMessage.Message)
			}
		}
	}

	if len(messages) > 0 {
		promptText = strings.Join(messages, "\n") + "\n" + promptText
	}

	return openai.ImageRequest{
		Model:          openai.CreateImageModelDallE3,
		Prompt:         promptText,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
}

func getMessage(update tg.Update, imgBytes []byte, user *models.User) (tg.Chattable, error) {
	file := createFile("image.jpg", imgBytes)
	err := messageService.CreateMessage(user.Id, update.Message.Text, "image.jpg", user.Mode)
	if err != nil {
		return nil, err
	}

	userMode := user.Mode

	if user.Dialog == 1 {
		userMode = userMode + 100
	}

	err = usageService.SaveUsage(user, userMode)
	if err != nil {
		return nil, fmt.Errorf("error saving usage: %w", err)
	}

	message := createMessage(update, file)

	return message, nil
}

func createFile(fileName string, imgBytes []byte) tg.FileBytes {
	return tg.FileBytes{
		Name:  fileName,
		Bytes: imgBytes,
	}
}

func createMessage(update tg.Update, file tg.FileBytes) tg.Chattable {
	message := tg.NewPhoto(update.Message.Chat.ID, file)
	message.ReplyToMessageID = update.Message.MessageID

	return message
}
