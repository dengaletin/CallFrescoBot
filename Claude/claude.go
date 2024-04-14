package claude

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
	anthropic "github.com/liushuangls/go-anthropic/v2"
	"log"
)

func handleResponse(update tg.Update, user *models.User, response anthropic.MessagesResponse) (tg.Chattable, error) {
	err := messageService.CreateMessage(user.Id, update.Message.Text, *response.Content[0].Text, user.Mode)
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

	text := response.Content[0].Text

	message := tg.NewMessage(update.Message.Chat.ID, *text)
	message.ReplyToMessageID = update.Message.MessageID

	return message, nil
}

func GetResponse(update tg.Update, user *models.User) (tg.Chattable, error) {
	if utils.GetEnvVar("CLAUDE_API_KEY") == "" {
		return nil, errors.New(consts.ErrorMissingClaudeKey)
	}

	client := anthropic.NewClient(utils.GetEnvVar("CLAUDE_API_KEY"))

	sendMsgErr := messageService.SendMsgToUser(update.Message.Chat.ID, utils.LocalizeSafe(consts.ClaudeLoading))
	if sendMsgErr != nil {
		return nil, sendMsgErr
	}

	var messages []anthropic.Message

	if user.Dialog == 1 {
		userMessages, err := messageService.GetMessagesByUser(user, 15, user.Mode)
		if err != nil {
			log.Printf("error getting messages by user: %v", err)
		}

		for _, userMessage := range userMessages {
			messages = append(messages, anthropic.NewUserTextMessage(userMessage.Message))
			if userMessage.Response != "" {
				messages = append(messages, anthropic.NewAssistantTextMessage(userMessage.Response))
			}
		}
	}

	messages = append(messages, anthropic.NewUserTextMessage(update.Message.Text))

	resp, err := client.CreateMessages(context.Background(), anthropic.MessagesRequest{
		Model:     anthropic.ModelClaudeInstant1Dot2,
		Messages:  messages,
		MaxTokens: 1000,
	})

	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
		return nil, err
	}

	msg, err := handleResponse(update, user, resp)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
