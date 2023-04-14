package commands

import "CallFrescoBot/pkg/models"

type Factory interface {
	RunCommand() ICommand
}

func GetCommand(cmd string, user *models.User) ICommand {
	switch cmd {
	default:
		return GptCommand{Message: cmd, User: user}
	case Start:
		return StartCommand{}
	case Status:
		return StatusCommand{User: user}
	}
}
