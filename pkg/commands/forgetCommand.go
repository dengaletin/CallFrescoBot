package commands

import (
	"CallFrescoBot/pkg/consts"
	userService "CallFrescoBot/pkg/service/user"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ForgetCommand struct {
	BaseCommand
}

func (cmd ForgetCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	err = userService.SetUserDialogFromId(cmd.User)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg), err
	}

	return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.Forget), nil
}
