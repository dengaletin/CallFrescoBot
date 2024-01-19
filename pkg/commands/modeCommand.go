package commands

import (
	"CallFrescoBot/pkg/consts"
	subsciptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type ModeCommand struct {
	BaseCommand
}

func (cmd ModeCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	modeID, idErr := getModeIDFromCommand(cmd.Update.Message.Text)
	if idErr != nil {
		return createErrorMessage(cmd.Update.Message.Chat.ID), idErr
	}

	subscription, subErr := subsciptionService.GetUserSubscription(cmd.User)
	if subErr != nil {
		return createErrorMessage(cmd.Update.Message.Chat.ID), subErr
	}

	if modeID == 2 && subscription == nil {
		return createCustomErrorMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.OnlyForPremium)),
			errors.New("only for premium users")
	}

	updateErr := userService.SetMode(modeID, cmd.User)
	if updateErr != nil {
		return createErrorMessage(cmd.Update.Message.Chat.ID), updateErr
	}

	modeName, modeErr := userService.GetMode(modeID)
	if modeErr != nil {
		return createErrorMessage(cmd.Update.Message.Chat.ID), modeErr
	}

	return createSuccessMessage(cmd.Update.Message.Chat.ID, modeName), nil
}

func getModeIDFromCommand(commandText string) (int64, error) {
	mode := strings.TrimPrefix(commandText, "/mode")
	return strconv.ParseInt(mode, 10, 64)
}

func createErrorMessage(chatID int64) tg.MessageConfig {
	return tg.NewMessage(chatID, utils.LocalizeSafe(consts.ErrorMsg))
}

func createCustomErrorMessage(chatID int64, text string) tg.MessageConfig {
	return tg.NewMessage(chatID, text)
}

func createSuccessMessage(chatID int64, modeName string) tg.MessageConfig {
	return tg.NewMessage(chatID, fmt.Sprintf(utils.LocalizeSafe(consts.ModeSuccess), modeName))
}
