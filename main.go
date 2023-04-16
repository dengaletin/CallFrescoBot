package main

import (
	"CallFrescoBot/app"
	"CallFrescoBot/pkg/commands"
	UserService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	app.SetupApp()

	upd := tg.NewUpdate(0)
	upd.Timeout = 30

	bot := utils.GetBot()
	updates := bot.GetUpdatesChan(upd)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			user, err := UserService.GetOrCreate(update.Message.From)
			if err != nil {
				log.Printf(err.Error())
				continue
			}

			response := commands.GetCommand(update.Message.Text, user).RunCommand()

			SendMessage(update, bot, response)
		}
	}
}

func SendMessage(update tg.Update, bot *tg.BotAPI, msgText string) {
	message := tg.NewMessage(update.Message.Chat.ID, msgText)
	message.ReplyToMessageID = update.Message.MessageID

	_, err := bot.Send(message)
	if err != nil {
		log.Printf(err.Error())
	}
}
