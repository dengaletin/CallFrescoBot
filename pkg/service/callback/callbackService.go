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
	Extra string `json:"extra"`
}

func ResolveAndHandle(query *tgbotapi.CallbackQuery, user *models.User, bot *tgbotapi.BotAPI) error {
	var data QueryData
	if err := json.Unmarshal([]byte(query.Data), &data); err != nil {
		log.Printf("Error while parsing query data: %s", err)
		return err
	}

	switch data.Type {
	case "mode":
		err := handleMode(data, user, bot, query)
		if err != nil {
			return err
		}
	case "context":
		err := handleContext(data, user, bot, query)
		if err != nil {
			return err
		}
	case "open":
		err := handleOpen(data, user, bot, query)
		if err != nil {
			return err
		}
	case "language":
		err := handleLanguage(data, user, bot, query)
		if err != nil {
			return err
		}
	default:
		log.Printf("Unknown query type: %s", data.Type)

		return errors.New("unknown query type")
	}

	return nil
}

func handleOpen(data QueryData, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	keyboard := "main"

	switch data.Value {
	case "main":
		keyboard = "main"
	case "language":
		keyboard = "language"
	case "buyLink":
		keyboard = "buyLink"
	case "buy":
		keyboard = "buy"
	default:
		log.Printf("Unknown open type: %s", data.Value)
	}

	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	nk, err := numericKeyboard.CreateNumericKeyboard(keyboard, user, data.Extra)
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

func handleLanguage(data QueryData, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	language, err := strconv.ParseInt(data.Value, 10, 64)
	if err != nil {
		log.Printf("handleMode: parse int64 error: %v", err)
		return err
	}

	err = userService.SetLanguage(language, user)
	if err != nil {
		return err
	}

	utils.InitBundle(user.Lang)

	mainMenuErr := numericKeyboard.CreateMainMenu()
	if mainMenuErr != nil {
		log.Printf(mainMenuErr.Error())
	}

	nk, err := numericKeyboard.CreateNumericKeyboard("language", user, data.Extra)
	if err != nil {
		log.Printf("Error creating keyboard: %v", err)
		return err
	}

	if nk == nil {
		log.Printf("Error: numeric keyboard is nil")

		return errors.New("numeric keyboard is nil")
	}

	message := utils.LocalizeSafe(consts.StartMsg)
	if data.Extra == "options" {
		message = utils.LocalizeSafe(consts.OptionsMessage)
	}

	_, err = bot.Send(tgbotapi.NewEditMessageText(
		query.Message.Chat.ID,
		query.Message.MessageID,
		message,
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

func handleMode(data QueryData, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	mode, err := strconv.ParseInt(data.Value, 10, 64)
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

	nk, err := numericKeyboard.CreateNumericKeyboard(data.Extra, user, data.Extra)
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

func handleContext(data QueryData, user *models.User, bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) error {
	_, err := bot.Request(tgbotapi.NewCallback(query.ID, ""))
	if err != nil {
		log.Printf("Error while responding to callback query: %s", err)
		return err
	}

	context, err := strconv.ParseInt(data.Value, 10, 64)
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

	nk, err := numericKeyboard.CreateNumericKeyboard(data.Extra, user, data.Extra)
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
