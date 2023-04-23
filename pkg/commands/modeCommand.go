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

type ModeCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd ModeCommand) Common() (string, error) {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}

func (cmd ModeCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common()
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	mode := strings.TrimPrefix(cmd.Update.Message.Text, "/mode")
	id, err := strconv.ParseInt(mode, 10, 64)
	if err != nil {
		return nil, err
	}

	err = userService.SetMode(id, cmd.User)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg), err
	}

	modeName, err := userService.GetMode(id)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg), err
	}

	return tg.NewMessage(cmd.Update.Message.Chat.ID, fmt.Sprintf(consts.ModeSuccess, modeName)), nil
}
