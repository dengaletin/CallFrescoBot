package numericKeyboard

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	payService "CallFrescoBot/pkg/service/invoice"
	planService "CallFrescoBot/pkg/service/plan"
	"CallFrescoBot/pkg/types"
	"CallFrescoBot/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
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
	Extra string `json:"extra"`
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
			Command:     "/status",
			Description: utils.LocalizeSafe(consts.StatusCommandDescription),
		},
		tg.BotCommand{
			Command:     "/buy",
			Description: utils.LocalizeSafe(consts.BuyCommandDescription),
		},
		//tg.BotCommand{
		//	Command:     "/invite",
		//	Description: utils.LocalizeSafe(consts.InviteCommandDescription),
		//},
		tg.BotCommand{
			Command:     "/options",
			Description: utils.LocalizeSafe(consts.OptionsCommandDescription),
		},
	)

	_, err := bot.Request(cmdCfg)
	if err != nil {
		return err
	}

	return nil
}

func CreateNumericKeyboard(keyboardType string, user *models.User, extra string) (*tg.InlineKeyboardMarkup, error) {
	switch keyboardType {
	case "main":
		return createMainKeyboard(user, extra), nil
	case "options":
		return createMainKeyboard(user, extra), nil
	case "language":
		return createLanguageKeyboard(user, extra), nil
	case "buy":
		return createBuyKeyboard(user, extra), nil
	case "buyLink":
		return createBuyLinkKeyboard(user, extra), nil
	default:
		fmt.Print(keyboardType)
		return nil, errors.New("unknown keyboard type")
	}
}

func createLanguageKeyboard(user *models.User, extra string) *tg.InlineKeyboardMarkup {
	englishButton := createButtonWithLang(utils.LocalizeSafe(consts.EnglishLanguage), user.Lang, 1, extra)
	russianButton := createButtonWithLang(utils.LocalizeSafe(consts.RussianLanguage), user.Lang, 2, extra)
	backButton := createButtonBack(utils.LocalizeSafe(consts.BackButton), extra)

	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(englishButton, russianButton),
		tg.NewInlineKeyboardRow(backButton),
	)

	return &keyboard
}

func createMainKeyboard(user *models.User, extra string) *tg.InlineKeyboardMarkup {
	chatGPTButton := createButtonWithMode("GPT3.5", extra, user.Mode, ModeChatGPT35)
	dalleButton := createButtonWithMode("DallE3", extra, user.Mode, ModeDallE)
	chatGPT4Button := createButtonWithMode("GPT4", extra, user.Mode, ModeChatGPT4)
	contextButton := createButtonWithContext(utils.LocalizeSafe(consts.ContextSupportButton), extra, user.Dialog)
	languageButton := createButtonWithLanguage(utils.LocalizeSafe(consts.LanguageSelectButton), extra)

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

func createButtonWithMode(text string, extra string, mode int64, buttonMode int64) tg.InlineKeyboardButton {
	activePrefix := "✅ "
	if mode == buttonMode {
		text = activePrefix + text
	}
	payload := createPayloadData(keyboardPayload{
		Type:  "mode",
		Value: strconv.FormatInt(buttonMode, 10),
		Extra: extra,
	})
	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonWithLang(text string, mode int64, buttonMode int64, extra string) tg.InlineKeyboardButton {
	activePrefix := "✅ "
	if mode == buttonMode {
		text = activePrefix + text
	}
	payload := createPayloadData(keyboardPayload{
		Type:  "language",
		Value: strconv.FormatInt(buttonMode, 10),
		Extra: extra,
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonWithContext(text string, extra string, dialog int64) tg.InlineKeyboardButton {
	activePrefix := "✅ "
	buttonValue := DialogOn
	if dialog == DialogOn {
		buttonValue = DialogOff
		text = activePrefix + text
	}
	payload := createPayloadData(keyboardPayload{
		Type:  "context",
		Value: strconv.Itoa(buttonValue),
		Extra: extra,
	})
	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonWithLanguage(text string, extra string) tg.InlineKeyboardButton {
	payload := createPayloadData(keyboardPayload{
		Type:  "open",
		Value: "language",
		Extra: extra,
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createButtonBack(text string, extra string) tg.InlineKeyboardButton {
	payload := createPayloadData(keyboardPayload{
		Type:  "open",
		Value: "main",
		Extra: extra,
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createBuyKeyboard(user *models.User, extra string) *tg.InlineKeyboardMarkup {
	plans, err := planService.GetAllPlans()
	if err != nil {
		log.Println("Failed to fetch plans:", err)
		return nil
	}

	var rows [][]tg.InlineKeyboardButton  // Сохраняем ряды кнопок
	var tempRow []tg.InlineKeyboardButton // Временная переменная для единого ряда

	for i, p := range plans {
		var config types.Config
		err := json.Unmarshal(p.Config, &config)
		if err != nil {
			log.Println("Failed to unmarshal plan config:", err)
			continue
		}

		currencySign := "$"
		planPrice := config.PriceEn
		if user.Lang == consts.LangRu {
			planPrice = config.PriceRu
			currencySign = "₽"
		}

		planName := p.Name
		buttonText := planName + " - " + currencySign + strconv.FormatFloat(planPrice, 'f', 2, 64)

		button := createBuyButton(buttonText, strconv.FormatUint(p.Id, 10))

		tempRow = append(tempRow, button)      // Добавляем кнопку в временный ряд
		if (i+1)%2 == 0 || i == len(plans)-1 { // Проверяем, если у нас две кнопки, или если это последний элемент
			rows = append(rows, tg.NewInlineKeyboardRow(tempRow...)) // Добавляем временный ряд в rows
			tempRow = []tg.InlineKeyboardButton{}                    // Сбрасываем временный ряд для следующих кнопок
		}
	}

	keyboard := tg.NewInlineKeyboardMarkup(rows...) // Создаем клавиатуру из всех рядов

	return &keyboard
}

func createBuyButton(text string, extra string) tg.InlineKeyboardButton {
	payload := createPayloadData(keyboardPayload{
		Type:  "open",
		Value: "buyLink",
		Extra: extra,
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}

func createBuyLinkKeyboard(user *models.User, extra string) *tg.InlineKeyboardMarkup {
	planID, err := strconv.ParseUint(extra, 10, 64)
	if err != nil {
		log.Println("Failed parse plan", err)
		return nil
	}

	plan, err := planService.GetPlanById(planID)
	if err != nil {
		log.Println("Failed get plan", err)
		return nil
	}

	url, err := payService.CreateInvoiceUrl(plan, user)
	if err != nil {
		log.Println("Failed to create invoice:", err)

		return nil
	}

	buttonText := utils.LocalizeSafe(consts.Buy) + " " + plan.Name

	urlButton := tg.NewInlineKeyboardButtonURL(buttonText, url)
	backButton := createButtonBuyBack(utils.LocalizeSafe(consts.BackButton), "")

	keyboard := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(urlButton),
		tg.NewInlineKeyboardRow(backButton),
	)

	return &keyboard
}

func createButtonBuyBack(text string, extra string) tg.InlineKeyboardButton {
	payload := createPayloadData(keyboardPayload{
		Type:  "open",
		Value: "buy",
		Extra: extra,
	})

	return tg.NewInlineKeyboardButtonData(text, payload)
}
