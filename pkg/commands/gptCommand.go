package commands

import (
	gpt "CallFrescoBot/Gpt"
	"CallFrescoBot/pkg/consts"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GptCommand struct {
	BaseCommand
}

func (cmd GptCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(true)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	gptResponse, err := gpt.GetResponse(cmd.Update, cmd.User)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg), err
	}

	return gptResponse, nil
}
