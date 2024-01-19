package callbackService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/service/numericKeyboard"
	subsciptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	"encoding/json"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

type QueryData struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func ResolveAndHandle(query *tgbotapi.CallbackQuery, user *models.User, bot *tgbotapi.BotAPI) error {
	var data QueryData
	if err := json.Unmarshal([]byte(query.Data), &data); err != nil {
		log.Printf("Error while parsing query data: %s", err)
		return err
	}

	switch data.Type {
	case "mode":
		err := handleMode(data.Value, user, bot, query)
		if err != nil {
			return err
		}
	case "context":
		err := handleContext(data.Value, user, bot, query)
		if err != nil {
			return err
		}
	case "open":
		err := handleOpen(data.Value, user, bot, query)
		if err != nil {
			return err
		}
	case "language":
		err := handleLanguage(data.Value, user, bot, query)
		if err != nil {
			return err
		}
	default:
		log.Printf("Unknown query type: %s", data.Type)

		return errors.New("unknown query type")
	}

	return nil
}

func handleOpen(value string, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	keyboard := "settings"

	switch value {
	case "settings":
		keyboard = "settings"
	case "language":
		keyboard = "language"
	default:
		log.Printf("Unknown open type: %s", value)
	}

	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	nk, err := numericKeyboard.CreateNumericKeyboard(keyboard, user)
	if err != nil {
		log.Printf("Error creating keyboard: %v", err)
		return err
	}

	_, err = bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		*nk,
	))
	if err != nil {
		return err
	}

	return nil
}

func handleLanguage(value string, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	language, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("handleMode: parse int64 error: %v", err)
		return err
	}

	err = userService.SetLanguage(language, user)
	if err != nil {
		return err
	}

	locErr := utils.CreateLoc(user)
	if locErr != nil {
		log.Printf(locErr.Error())
	}

	mainMenuErr := numericKeyboard.CreateMainMenu()
	if mainMenuErr != nil {
		log.Printf(mainMenuErr.Error())
	}

	nk, err := numericKeyboard.CreateNumericKeyboard("language", user)
	if err != nil {
		log.Printf("Error creating keyboard: %v", err)
		return err
	}

	if nk == nil {
		log.Printf("Error: numeric keyboard is nil")

		return errors.New("numeric keyboard is nil")
	}

	_, err = bot.Send(tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		utils.LocalizeSafe(consts.SettingsMessage),
	))
	if err != nil {
		return err
	}

	_, err = bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		*nk,
	))
	if err != nil {
		return err
	}

	return nil
}

func handleMode(value string, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	mode, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("handleMode: parse int64 error: %v", err)
		return err
	}

	subscription, subErr := subsciptionService.GetUserSubscription(user)
	if subErr != nil {
		return subErr
	}

	if mode == 2 && subscription == nil {
		_, sendErr := bot.Send(tgbotapi.NewMessage(query.Message.Chat.ID, utils.LocalizeSafe(consts.OnlyForPremium)))
		if sendErr != nil {
			return sendErr
		}
		return errors.New("only for premium users")
	}

	err = userService.SetMode(mode, user)
	if err != nil {
		return err
	}

	nk, err := numericKeyboard.CreateNumericKeyboard("settings", user)
	if err != nil {
		log.Printf("Error creating keyboard: %v", err)
		return err
	}

	if nk == nil {
		log.Printf("Error: numeric keyboard is nil")
		return errors.New("numeric keyboard is nil")
	}

	_, err = bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		*nk,
	))
	if err != nil {
		return err
	}

	return nil
}

func handleContext(value string, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	context, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("handleContext: parse int64 error: %v", err)
		return err
	}

	subscription, subErr := subsciptionService.GetUserSubscription(user)
	if subErr != nil {
		return subErr
	}

	if context == 1 && subscription == nil {
		_, sendErr := bot.Send(tgbotapi.NewMessage(query.Message.Chat.ID, utils.LocalizeSafe(consts.OnlyForPremium)))
		if sendErr != nil {
			return sendErr
		}
		return errors.New("only for premium users")
	}

	err = userService.SetDialogStatus(context, user)
	if err != nil {
		return err
	}

	nk, err := numericKeyboard.CreateNumericKeyboard("settings", user)
	if err != nil {
		log.Printf("Error creating keyboard: %v", err)
		return err
	}

	_, err = bot.Send(tgbotapi.NewEditMessageReplyMarkup(
		query.Message.Chat.ID,
		query.Message.MessageID,
		*nk,
	))

	if err != nil {
		return err
	}

	return nil
}
