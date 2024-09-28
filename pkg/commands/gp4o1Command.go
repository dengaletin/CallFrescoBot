package commands

import (
	gpt "CallFrescoBot/Gpt"
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

type Gpt4o1Command struct {
	BaseCommand
}

func (cmd Gpt4o1Command) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(true)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	response, err := gpt.GetResponse(cmd.Update, cmd.User, openai.O1Preview)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.ErrorMsg))}, err
	}

	return response, nil
}
