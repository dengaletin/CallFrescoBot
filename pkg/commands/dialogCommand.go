package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	userService "CallFrescoBot/pkg/service/user"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type DialogCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd DialogCommand) Common() (string, error) {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}

func (cmd DialogCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common()
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	dialog := strings.TrimPrefix(cmd.Update.Message.Text, "/dialog")
	id, err := strconv.ParseInt(dialog, 10, 64)
	if err != nil {
		return nil, err
	}

	err = userService.SetDialogStatus(id, cmd.User)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg), err
	}

	dialogStatus, err := userService.GetDialogStatus(id)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg), err
	}

	return tg.NewMessage(cmd.Update.Message.Chat.ID, fmt.Sprintf(consts.DialogSuccess, dialogStatus)), nil
}
