package commands

import (
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	subsciptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd *BaseCommand) Common(validateUser bool) (string, error) {
	if validateUser {
		userValidatorMessage, err := userService.ValidateUser(cmd.User)
		if err != nil {
			return userValidatorMessage, err
		}
	}

	err := subsciptionService.ResetSubscription(cmd.User)
	if err != nil {
		return err.Error(), err
	}

	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}
