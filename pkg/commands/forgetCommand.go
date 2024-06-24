package commands

import (
	"CallFrescoBot/pkg/consts"
	userService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ForgetCommand struct {
	BaseCommand
}

func (cmd ForgetCommand) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	err = userService.SetUserDialogFromId(cmd.User)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.ErrorMsg))}, err
	}

	return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.Forget))}, nil
}
