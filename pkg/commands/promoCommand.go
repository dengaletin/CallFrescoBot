package commands

import (
	"CallFrescoBot/pkg/consts"
	campaignService "CallFrescoBot/pkg/service/campaign"
	"CallFrescoBot/pkg/service/numericKeyboard"
	promoService "CallFrescoBot/pkg/service/promo"
	"CallFrescoBot/pkg/utils"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type PromoCommand struct {
	BaseCommand
}

func (cmd PromoCommand) RunCommand() ([]tg.Chattable, error) {
	msg := tg.NewMessage(cmd.Update.Message.Chat.ID, "")

	nk, err := numericKeyboard.CreateNumericKeyboard("main", cmd.User, "main")
	if err != nil {
		msg.Text = "Ошибка создания клавиатуры"
	}

	msg.ReplyMarkup = nk
	msg.ParseMode = "markdown"

	result, err := cmd.Common(false)
	if err != nil {
		msg.Text = result
		return []tg.Chattable{msg}, err
	}

	promoCode, err := parsePromoCode(cmd.Update.Message.Text)
	if err != nil {
		msg.Text = utils.LocalizeSafe(consts.StartMsg)
		return []tg.Chattable{msg}, err
	}

	if !cmd.User.IsNew {
		msg.Text = utils.LocalizeSafe(consts.StartMsg)
		return []tg.Chattable{msg}, errors.New("user is not new")
	}

	campaign, err := campaignService.Get(promoCode)
	if err != nil {
		msg.Text = utils.LocalizeSafe(consts.StartMsg)
		return []tg.Chattable{msg}, err
	}

	if err = promoService.Create(campaign.Id, cmd.User); err != nil {
		msg.Text = utils.LocalizeSafe(consts.StartMsg)
		return []tg.Chattable{msg}, err
	}

	msg.Text = utils.LocalizeSafe(consts.StartMsg)

	return []tg.Chattable{msg}, err
}

func parsePromoCode(messageText string) (string, error) {
	parts := strings.SplitN(messageText, " ", 2)

	if len(parts) != 2 {
		return "", errors.New("invalid message format")
	}

	promoCode := parts[1]

	return promoCode, nil
}
