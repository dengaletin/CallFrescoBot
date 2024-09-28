package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
)

const (
	StartCommandPattern = `^/start ref[0-9]+$`
	PromoCommandPattern = `^/start [A-Za-z]+$`
)

type CommandRegistryEntry struct {
	Pattern   *regexp.Regexp
	Generator func(update tgbotapi.Update, user *models.User) ICommandInterface
}

var commandRegistry = []CommandRegistryEntry{
	{Pattern: regexp.MustCompile(StartCommandPattern), Generator: NewRefCommand},
	{Pattern: regexp.MustCompile(PromoCommandPattern), Generator: NewPromoCommand},
}

func NewStartCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	if user.IsNew {
		return FirstRun{BaseCommand{Update: update, User: user}}
	}
	return StartCommand{BaseCommand{Update: update, User: user}}
}

func NewRefCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return RefCommand{BaseCommand{Update: update, User: user}}
}

func NewStatusCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return StatusCommand{BaseCommand{Update: update, User: user}}
}

func NewOptionsCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return OptionsCommand{BaseCommand{Update: update, User: user}}
}

func NewInviteCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return InviteCommand{BaseCommand{Update: update, User: user}}
}

func NewPromoCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return PromoCommand{BaseCommand{Update: update, User: user}}
}

func NewBuyCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return BuyCommand{BaseCommand{Update: update, User: user}}
}

func NewForgetCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return ForgetCommand{BaseCommand{Update: update, User: user}}
}

func NewDalleCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return DalleCommand{BaseCommand{Update: update, User: user}}
}

func NewGptCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
	return GptCommand{BaseCommand{Update: update, User: user}}
}

func NewGpt4Command(update tgbotapi.Update, user *models.User) ICommandInterface {
	return Gpt4Command{BaseCommand{Update: update, User: user}, utils.GetBot()}
}

func NewGpt4o1Command(update tgbotapi.Update, user *models.User) ICommandInterface {
	return Gpt4o1Command{BaseCommand{Update: update, User: user}}
}

func GetCommand(update tgbotapi.Update, user *models.User) ICommandInterface {
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
	case "/buy":
		return NewBuyCommand(update, user)
	case "/forget":
		return NewForgetCommand(update, user)
	case "/options":
		return NewOptionsCommand(update, user)
	default:
		if user.Mode == consts.DalleMode {
			return NewDalleCommand(update, user)
		}
		if user.Mode == consts.Gpt4oMode {
			return NewGpt4Command(update, user)
		}
		if user.Mode == consts.Gpt4o1Mode {
			return NewGpt4o1Command(update, user)
		}
		return NewGptCommand(update, user)
	}
}
