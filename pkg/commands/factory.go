package commands

import (
	"CallFrescoBot/pkg/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

const (
	StartCommandPattern = `^/start ref[0-9]+$`
	ModeCommandPattern  = `^/mode[0-9]$`
)

type CommandRegistryEntry struct {
	Pattern   *regexp.Regexp
	Generator func(update tgbotapi.Update, user *models.User) ICommand
}

var commandRegistry = []CommandRegistryEntry{
	{Pattern: regexp.MustCompile(StartCommandPattern), Generator: NewRefCommand},
	{Pattern: regexp.MustCompile(ModeCommandPattern), Generator: NewModeCommand},
}

func NewModeCommand(update tgbotapi.Update, user *models.User) ICommand {
	return ModeCommand{BaseCommand{Update: update, User: user}}
}

func NewStartCommand(update tgbotapi.Update, user *models.User) ICommand {
	return StartCommand{BaseCommand{Update: update, User: user}}
}

func NewRefCommand(update tgbotapi.Update, user *models.User) ICommand {
	return RefCommand{BaseCommand{Update: update, User: user}}
}

func NewStatusCommand(update tgbotapi.Update, user *models.User) ICommand {
	return StatusCommand{BaseCommand{Update: update, User: user}}
}

func NewOptionsCommand(update tgbotapi.Update, user *models.User) ICommand {
	return OptionsCommand{BaseCommand{Update: update, User: user}}
}

func NewInviteCommand(update tgbotapi.Update, user *models.User) ICommand {
	return InviteCommand{BaseCommand{Update: update, User: user}}
}

func NewBuyCommand(update tgbotapi.Update, user *models.User) ICommand {
	return BuyCommand{BaseCommand{Update: update, User: user}}
}

func NewForgetCommand(update tgbotapi.Update, user *models.User) ICommand {
	return ForgetCommand{BaseCommand{Update: update, User: user}}
}

func NewDalleCommand(update tgbotapi.Update, user *models.User) ICommand {
	return DalleCommand{BaseCommand{Update: update, User: user}}
}

func NewGptCommand(update tgbotapi.Update, user *models.User) ICommand {
	return GptCommand{BaseCommand{Update: update, User: user}}
}

func NewGpt4Command(update tgbotapi.Update, user *models.User) ICommand {
	return Gpt4Command{BaseCommand{Update: update, User: user}}
}

func GetCommand(update tgbotapi.Update, user *models.User) ICommand {
	for _, entry := range commandRegistry {
		if entry.Pattern.MatchString(update.Message.Text) {
			return entry.Generator(update, user)
		}
	}

	switch update.Message.Text {
	case "/start":
		return NewStartCommand(update, user)
	case "/status":
		return NewStatusCommand(update, user)
	case "/invite":
		return NewInviteCommand(update, user)
	case "/buy":
		return NewBuyCommand(update, user)
	case "/forget":
		return NewForgetCommand(update, user)
	case "/options":
		return NewOptionsCommand(update, user)
	default:
		if user.Mode == 1 {
			return NewDalleCommand(update, user)
		}
		if user.Mode == 2 {
			return NewGpt4Command(update, user)
		}
		return NewGptCommand(update, user)
	}
}
