package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	userRefService "CallFrescoBot/pkg/service/userRef"
	"log"
	"strconv"
	"strings"
)

type RefCommand struct {
	Message string
	User    *models.User
}

func (cmd RefCommand) Common() string {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Message)
	if err != nil {
		log.Printf(err.Error())
		return messageValidatorText
	}

	return ""
}

func (cmd RefCommand) RunCommand() string {
	result := cmd.Common()

	if result != "" {
		return result
	}

	message := strings.TrimPrefix(cmd.Message, "/start ref")

	id, err := strconv.ParseInt(message, 10, 64)
	if err != nil {
		log.Printf(err.Error())
		return consts.StartMsg
	}

	if id == cmd.User.TgId {
		log.Printf("The user is trying to invite himself.")
		return consts.StartMsg
	}

	if !cmd.User.IsNew {
		log.Printf("The user is not new.")
		return consts.StartMsg
	}

	user, err := userService.GerUserByTgId(id)
	if err != nil {
		log.Printf(err.Error())
		return consts.StartMsg
	}

	_, err = userRefService.Create(user, cmd.User)
	if err != nil {
		log.Printf(err.Error())
		return consts.StartMsg
	}

	_, err = subscriptionService.GetOrCreate(user, 50, consts.RefDaysMultiplier)
	if err != nil {
		log.Printf(err.Error())
		return consts.StartMsg
	}

	return consts.StartMsg
}
