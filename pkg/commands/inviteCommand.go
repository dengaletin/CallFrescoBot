package commands

import (
	"CallFrescoBot/pkg/consts"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InviteCommand struct {
	BaseCommand
}

func (cmd InviteCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	inviteLink := fmt.Sprintf(consts.InviteLink, cmd.User.TgId)

	return tg.NewMessage(cmd.Update.Message.Chat.ID, inviteLink), nil
}
