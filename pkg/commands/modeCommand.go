package commands

import (
	"CallFrescoBot/pkg/consts"
	userService "CallFrescoBot/pkg/service/user"
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
	return tg.NewMessage(chatID, consts.ErrorMsg)
}

func createSuccessMessage(chatID int64, modeName string) tg.MessageConfig {
	return tg.NewMessage(chatID, fmt.Sprintf(consts.ModeSuccess, modeName))
}
