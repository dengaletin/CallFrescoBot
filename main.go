package main

import (
	"CallFrescoBot/pkg/commands"
	"CallFrescoBot/pkg/messages"
	"CallFrescoBot/pkg/validator"
	"log"
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	apiKey := os.Getenv("TELEGRAM_API_KEY")
	if apiKey == "" {
		log.Fatalln(messages.MissingTgKey)
	}

	bot, err := tg.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	upd := tg.NewUpdate(0)
	upd.Timeout = 30

	updates := bot.GetUpdatesChan(upd)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			messageText, err := validator.Validate(update.Message.Text)

			if err != nil {
				log.Printf(err.Error())
				SendMessage(update, bot, messageText)
				continue
			}

			response := commands.GetCommand(messageText).RunCommand()
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
