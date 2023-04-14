package commands

import (
	gpt "CallFrescoBot/Gpt"
	"CallFrescoBot/pkg/models"
)

type GptCommand struct {
	Message string
	User    *models.User
}

func (cmd GptCommand) RunCommand() string {
	return gpt.GetResponse(cmd.Message, cmd.User)
}
