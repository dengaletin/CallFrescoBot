package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/service/numericKeyboard"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type OptionsCommand struct {
	BaseCommand
}

func (cmd OptionsCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	nk, err := numericKeyboard.CreateNumericKeyboard("main", cmd.User, "options")

	msg := tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.OptionsMessage))
	msg.ReplyMarkup = nk

	return msg, nil
}
