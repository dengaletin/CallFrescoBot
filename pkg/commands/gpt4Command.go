package commands

import (
	gpt "CallFrescoBot/Gpt"
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

type Gpt4Command struct {
	BaseCommand
	Bot *tg.BotAPI
}

func (cmd Gpt4Command) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(true)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	gptResponses, err := gpt.GetResponse(cmd.Update, cmd.User, openai.GPT4o)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.ErrorMsg))}, err
	}

	return gptResponses, nil
}
