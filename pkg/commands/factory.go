package commands

import (
	"CallFrescoBot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

type Factory interface {
	RunCommand() ICommand
}

func GetCommand(cmd tgbotapi.Update, user *models.User) ICommand {
	re := regexp.MustCompile(`^/start ref[0-9]+$`)
	match := re.FindStringSubmatch(cmd.Message.Text)
	if len(match) != 0 {
		return RefCommand{Update: cmd, User: user}
	}

	reMode := regexp.MustCompile(`^/mode[0-9]$`)
	matchMode := reMode.FindStringSubmatch(cmd.Message.Text)
	if len(matchMode) != 0 {
		return ModeCommand{Update: cmd, User: user}
	}

	reDialog := regexp.MustCompile(`^/dialog[0-9]$`)
	matchDialog := reDialog.FindStringSubmatch(cmd.Message.Text)
	if len(matchDialog) != 0 {
		return DialogCommand{Update: cmd, User: user}
	}

	switch cmd.Message.Text {
	default:
		if user.Mode != 0 {
			return DalleCommand{Update: cmd, User: user}
		}

		return GptCommand{Update: cmd, User: user}
	case Start:
		return StartCommand{Update: cmd, User: user}
	case Status:
		return StatusCommand{Update: cmd, User: user}
	case Invite:
		return InviteCommand{Update: cmd, User: user}
	case Buy:
		return BuyCommand{Update: cmd, User: user}
	}
}
