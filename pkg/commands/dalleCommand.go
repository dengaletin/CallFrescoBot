package commands

import (
	"CallFrescoBot/Dalle3"
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DalleCommand struct {
	BaseCommand
}

func (cmd DalleCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(true)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	dalleResponse, err := Dalle3.GetResponse(cmd.Update, cmd.User)

	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.ErrorMsg)), err
	}

	return dalleResponse, nil
}
