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

var apiKey string

func init() {
	apiKey = utils.GetEnvVar("GPT_API_KEY")
}

func GetResponse(update tg.Update, user *models.User) (tg.Chattable, error) {
	if apiKey == "" {
		return nil, errors.New(consts.MissingGptKey)
	}

	c := openai.NewClient(apiKey)
	ctx := context.Background()

	reqBase64 := openai.ImageRequest{
		Model:          openai.CreateImageModelDallE3,
		Prompt:         update.Message.Text,
		Size:           openai.CreateImageSize1792x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		N:              1,
	}

	respBase64, err := c.CreateImage(ctx, reqBase64)
	if err != nil {
		fmt.Printf("Image creation error: %v\n", err)
		return nil, err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return nil, err
	}

	file := tg.FileBytes{
		Name:  "image.jpg",
		Bytes: imgBytes,
	}

	err = messageService.CreateMessage(user.Id, update.Message.Text, "image.jpg")
	if err != nil {
		return nil, err
	}

	message := tg.NewPhoto(update.Message.Chat.ID, file)
	message.ReplyToMessageID = update.Message.MessageID

	return message, nil
}
