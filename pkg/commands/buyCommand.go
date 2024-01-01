package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BuyCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd BuyCommand) Common() (string, error) {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}

func (cmd BuyCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common()
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.BuyMsg), nil
}
