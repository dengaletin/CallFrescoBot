package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InviteCommand struct {
	BaseCommand
}

func (cmd InviteCommand) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	inviteLink := fmt.Sprintf(utils.LocalizeSafe(consts.InviteLink), cmd.User.TgId)
	message := tg.NewMessage(cmd.Update.Message.Chat.ID, inviteLink)
	message.ParseMode = "markdown"

	return []tg.Chattable{message}, nil
}
