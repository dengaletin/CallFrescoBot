package commands

import (
	"CallFrescoBot/Dalle2"
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	userService "CallFrescoBot/pkg/service/user"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DalleCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd DalleCommand) Common() (string, error) {
	userValidatorMessage, err := userService.ValidateUser(cmd.User)
	if err != nil {
		return userValidatorMessage, err
	}

	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}

func (cmd DalleCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common()
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	dalleResponse, err := Dalle2.GetResponse(cmd.Update, cmd.User)

	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.ErrorMsg), err
	}

	return dalleResponse, nil
}
