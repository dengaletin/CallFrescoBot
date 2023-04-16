package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	"fmt"
	"log"
)

type InviteCommand struct {
	Message string
	User    *models.User
}

func (cmd InviteCommand) Common() string {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Message)
	if err != nil {
		log.Printf(err.Error())
		return messageValidatorText
	}

	return ""
}

func (cmd InviteCommand) RunCommand() string {
	result := cmd.Common()

	if result != "" {
		return result
	}

	inviteLink := fmt.Sprintf(consts.InviteLink, cmd.User.TgId)

	return inviteLink
}
