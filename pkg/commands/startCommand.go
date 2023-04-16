package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	"log"
)

type StartCommand struct {
	Message string
	User    *models.User
}

func (cmd StartCommand) Common() string {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Message)
	if err != nil {
		log.Printf(err.Error())
		return messageValidatorText
	}

	return ""
}

func (cmd StartCommand) RunCommand() string {
	result := cmd.Common()

	if result != "" {
		return result
	}

	return consts.StartMsg
}
