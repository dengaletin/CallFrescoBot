package commands

import (
	"CallFrescoBot/pkg/models"
	"regexp"
)

type Factory interface {
	RunCommand() ICommand
}

func GetCommand(cmd string, user *models.User) ICommand {
	re := regexp.MustCompile(`^/start ref[0-9]+$`)
	match := re.FindStringSubmatch(cmd)
	if len(match) != 0 {
		return RefCommand{Message: cmd, User: user}
	}

	switch cmd {
	default:
		return GptCommand{Message: cmd, User: user}
	case Start:
		return StartCommand{Message: cmd, User: user}
	case Status:
		return StatusCommand{Message: cmd, User: user}
	case Invite:
		return InviteCommand{Message: cmd, User: user}
	}
}
