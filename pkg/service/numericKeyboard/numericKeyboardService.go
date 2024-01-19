package numericKeyboard

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/utils"
	"encoding/json"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

const (
	ModeChatGPT35 = iota
	ModeDallE
	ModeChatGPT4
)

const (
	DialogOff = iota
	DialogOn
)

type keyboardPayload struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func CreateMainMenu() error {
	bot := utils.GetBot()

	cmdCfg := tg.NewSetMyCommands(
		tg.BotCommand{
			Command:     "/start",
			Description: utils.LocalizeSafe(consts.StartCommandDescription),
		},
		tg.BotCommand{
			Command:     "/forget",
			Description: utils.LocalizeSafe(consts.ForgetCommandDescription),
		},
		tg.BotCommand{
			Command:     "/buy",
			Description: utils.LocalizeSafe(consts.BuyCommandDescription),
		},
		tg.BotCommand{
			Command:     "/invite",
			Description: utils.LocalizeSafe(consts.InviteCommandDescription),
		},
		tg.BotCommand{
			Command:     "/settings",
			Description: utils.LocalizeSafe(consts.SettingsCommandDescription),
		},
	)

	_, err := bot.Request(cmdCfg)
	if err != nil {
		return err
	}

	return nil
}

func CreateNumericKeyboard(keyboardType string, user *models.User) (*tg.InlineKeyboardMarkup, error) {
	switch keyboardType {
	case "settings":
		return createSettingsKeyboard(user), nil
	case "language":
		return createLanguageKeyboard(user), nil
	default:
		return nil, errors.New("unknown keyboard type")
	}
}

func createLanguageKeyboard(user *models.User) *tg.InlineKeyboardMarkup {
	englishButton := createButtonWithLang(utils.LocalizeSafe(consts.EnglishLanguage), user.Lang, 1)
	russianButton := createButtonWithLang(utils.LocalizeSafe(consts.RussianLanguage), user.Lang, 2)
	backButton := createButtonBack(utils.LocalizeSafe(consts.BackButton))

	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(englishButton, russianButton),
		tg.NewInlineKeyboardRow(backButton),
	)

	return &keyboard
}

func createSettingsKeyboard(user *models.User) *tg.InlineKeyboardMarkup {
	chatGPTButton := createButtonWithMode("GPT3.5", user.Mode, ModeChatGPT35)
	dalleButton := createButtonWithMode("DallE3", user.Mode, ModeDallE)
	chatGPT4Button := createButtonWithMode("GPT4", user.Mode, ModeChatGPT4)
	contextButton := createButtonWithContext(utils.LocalizeSafe(consts.ContextSupportButton), user.Dialog)
	languageButton := createButtonWithLanguage(utils.LocalizeSafe(consts.LanguageSelectButton))

	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(chatGPTButton, dalleButton, chatGPT4Button),
		tg.NewInlineKeyboardRow(contextButton),
		tg.NewInlineKeyboardRow(languageButton),
	)

	return &keyboard
}

func createPayloadData(payload keyboardPayload) string {
	data, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(data)
}

func createButtonWithMode(text string, mode int64, buttonMode int64) tg.InlineKeyboardButton {
	activePrefix := "✅ "
	if mode == buttonMode {
		text = activePrefix + text
	}
	payload := createPayloadData(keyboardPayload{
		Type:  "mode",
		Value: strconv.FormatInt(buttonMode, 10),
	})
	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonWithLang(text string, mode int64, buttonMode int64) tg.InlineKeyboardButton {
	activePrefix := "✅ "
	if mode == buttonMode {
		text = activePrefix + text
	}
	payload := createPayloadData(keyboardPayload{
		Type:  "language",
		Value: strconv.FormatInt(buttonMode, 10),
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonWithContext(text string, dialog int64) tg.InlineKeyboardButton {
	activePrefix := "✅ "
	buttonValue := DialogOn
	if dialog == DialogOn {
		buttonValue = DialogOff
		text = activePrefix + text
	}
	payload := createPayloadData(keyboardPayload{
		Type:  "context",
		Value: strconv.Itoa(buttonValue),
	})
	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonWithLanguage(text string) tg.InlineKeyboardButton {
	payload := createPayloadData(keyboardPayload{
		Type:  "open",
		Value: "language",
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonBack(text string) tg.InlineKeyboardButton {
	payload := createPayloadData(keyboardPayload{
		Type:  "open",
		Value: "settings",
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}
