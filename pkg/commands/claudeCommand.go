package commands

import (
	claude "CallFrescoBot/Claude"
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ClaudeCommand struct {
	BaseCommand
}

func (cmd ClaudeCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(true)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	claudeResponse, err := claude.GetResponse(cmd.Update, cmd.User)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.ErrorMsg)), err
	}

	return claudeResponse, nil
}
