package commands

import (
	"CallFrescoBot/pkg/consts"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BuyCommand struct {
	BaseCommand
}

func (cmd BuyCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.BuyMsg), nil
}
