package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InviteCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd InviteCommand) Common() (string, error) {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}

func (cmd InviteCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common()
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	inviteLink := fmt.Sprintf(consts.InviteLink, cmd.User.TgId)

	return tg.NewMessage(cmd.Update.Message.Chat.ID, inviteLink), nil
}