package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/service/numericKeyboard"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FirstRun struct {
	BaseCommand
}

func (cmd FirstRun) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	nk, err := numericKeyboard.CreateNumericKeyboard("firstRun", cmd.User, "firstRun")

	msg := tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.FirstRunMsg))
	msg.ReplyMarkup = nk
	msg.ParseMode = "markdown"

	return []tg.Chattable{msg}, nil
}
