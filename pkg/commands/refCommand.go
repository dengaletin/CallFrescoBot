package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	userRefService "CallFrescoBot/pkg/service/userRef"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type RefCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd RefCommand) Common() (string, error) {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}

func (cmd RefCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common()
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	message := strings.TrimPrefix(cmd.Update.Message.Text, "/start ref")

	id, err := strconv.ParseInt(message, 10, 64)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.StartMsg), err
	}

	if id == cmd.User.TgId {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.StartMsg), errors.New("the user is trying to invite himself")
	}

	if !cmd.User.IsNew {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.StartMsg), errors.New("the user is not new")
	}

	user, err := userService.GerUserByTgId(id)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.StartMsg), err
	}

	_, err = userRefService.Create(user, cmd.User)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.StartMsg), err
	}

	_, err = subscriptionService.GetOrCreate(user, 25, consts.RefDaysMultiplier)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.StartMsg), err
	}

	return tg.NewMessage(cmd.Update.Message.Chat.ID, consts.StartMsg), nil
}
