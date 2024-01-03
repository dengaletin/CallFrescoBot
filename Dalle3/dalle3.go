package Dalle3

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	"CallFrescoBot/pkg/utils"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

var apiKey = utils.GetEnvVar("GPT_API_KEY")

func GetResponse(update tg.Update, user *models.User) (tg.Chattable, error) {
	if apiKey == "" {
		return nil, errors.New(consts.MissingGptKey)
	}

	return getImageResponse(update, user)
}

func getImageResponse(update tg.Update, user *models.User) (tg.Chattable, error) {
	c := openai.NewClient(apiKey)
	ctx := context.Background()
	imgBytes, err := getImageBytes(ctx, c, update)

	if err != nil {
		return nil, err
	}

	return getMessage(update, imgBytes, user.Id)
}

func getImageBytes(ctx context.Context, openaiClient *openai.Client, update tg.Update) ([]byte, error) {
	reqBase64 := prepareRequest(update.Message.Text)
	respBase64, err := openaiClient.CreateImage(ctx, reqBase64)

	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return nil, err
	}

	return base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
}

func prepareRequest(promptText string) openai.ImageRequest {
	return openai.ImageRequest{
		Model:          openai.CreateImageModelDallE3,
		Prompt:         promptText,
		Size:           openai.CreateImageSize1792x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}
}

func getMessage(update tg.Update, imgBytes []byte, userID uint64) (tg.Chattable, error) {
	file := createFile("image.jpg", imgBytes)
	err := messageService.CreateMessage(userID, update.Message.Text, "image.jpg")

	if err != nil {
		return nil, err
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
