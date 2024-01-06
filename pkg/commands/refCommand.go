package commands

import (
	"CallFrescoBot/pkg/consts"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	userRefService "CallFrescoBot/pkg/service/userRef"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type RefCommand struct {
	BaseCommand
}

func (cmd RefCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return newMessage(cmd.Update, result), err
	}

	refID, err := parseReferralID(cmd.Update.Message.Text)
	if err != nil {
		return newMessage(cmd.Update, consts.StartMsg), err
	}

	if refID == cmd.User.TgId {
		return newMessage(cmd.Update, consts.StartMsg), errors.New("cannot refer self")
	}

	if !cmd.User.IsNew {
		return newMessage(cmd.Update, consts.StartMsg), errors.New("user is not new")
	}

	referringUser, err := userService.GerUserByTgId(refID)
	if err != nil {
		return newMessage(cmd.Update, consts.StartMsg), err
	}

	if _, err = userRefService.Create(referringUser, cmd.User); err != nil {
		return newMessage(cmd.Update, consts.StartMsg), err
	}

	if _, err = subscriptionService.GetOrCreate(referringUser, 25, consts.RefDaysMultiplier); err != nil {
		return newMessage(cmd.Update, consts.StartMsg), err
	}

	return newMessage(cmd.Update, consts.StartMsg), nil
}

func newMessage(update tg.Update, text string) tg.Chattable {
	return tg.NewMessage(update.Message.Chat.ID, text)
}

func parseReferralID(messageText string) (int64, error) {
	message := strings.TrimPrefix(messageText, "/start ref")
	return strconv.ParseInt(strings.TrimSpace(message), 10, 64)
}
