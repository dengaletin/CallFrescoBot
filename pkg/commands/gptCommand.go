package commands

import (
	gpt "CallFrescoBot/Gpt"
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

type GptCommand struct {
	BaseCommand
}

func (cmd GptCommand) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(true)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	bot := utils.GetBot()
	response, err := gpt.GetResponse(bot, cmd.Update, cmd.User, openai.GPT4oMini)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.ErrorMsg))}, err
	}

	return response, nil
}
