package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/service/numericKeyboard"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
)

type BuyCommand struct {
	BaseCommand
}

func (cmd BuyCommand) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	nk, err := numericKeyboard.CreateNumericKeyboard("buy", cmd.User, "buy")

	msg := tg.NewMessage(cmd.Update.Message.Chat.ID, utils.LocalizeSafe(consts.BuyMsg))
	msg.ReplyMarkup = nk
	msg.ParseMode = "markdown"

	var filename string

	if cmd.User.Lang == consts.LangEn {
		filename = "subscriptions.jpg"
	} else {
		filename = "subscriptions_ru.jpg"
	}

	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	photoFileBytes := tg.FileBytes{
		Name:  "subscriptions.jpg",
		Bytes: bytes,
	}

	photoMsg := tg.NewPhoto(cmd.Update.Message.Chat.ID, photoFileBytes)
	photoMsg.Caption = msg.Text
	photoMsg.ParseMode = msg.ParseMode
	photoMsg.ReplyMarkup = msg.ReplyMarkup

	return []tg.Chattable{photoMsg}, nil
}
