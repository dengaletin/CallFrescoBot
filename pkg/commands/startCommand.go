package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartCommand struct {
	BaseCommand
}

func (cmd StartCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	return tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.StartMsg)), nil
}
