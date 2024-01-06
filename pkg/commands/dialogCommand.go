package commands

import (
	"CallFrescoBot/pkg/consts"
	userService "CallFrescoBot/pkg/service/user"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type DialogCommand struct {
	BaseCommand
}

func (cmd DialogCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	id, err := extractDialogID(cmd.Update.Message.Text)
	if err != nil {
		return nil, err
	}

	err = userService.SetDialogStatus(id, cmd.User)
	if err != nil {
		return cmd.newErrorMessage(), err
	}

	dialogStatus, err := userService.GetDialogStatus(id)
	if err != nil {
		return cmd.newErrorMessage(), err
	}

	successMsg := fmt.Sprintf(consts.DialogSuccess, dialogStatus)
	return tg.NewMessage(cmd.Update.Message.Chat.ID, successMsg), nil
}

func extractDialogID(messageText string) (int64, error) {
	dialog := strings.TrimPrefix(messageText, "/dialog")
	return strconv.ParseInt(dialog, 10, 64)
}

func (cmd DialogCommand) newErrorMessage() tg.Chattable {
	return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg)
}
