package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/service/numericKeyboard"
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

	nk, err := numericKeyboard.CreateNumericKeyboard("main", cmd.User, "main")

	msg := tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.StartMsg))
	msg.ReplyMarkup = nk
	msg.ParseMode = "markdown"
	return msg, nil
}
