package main

import (
	"CallFrescoBot/pkg/commands"
	"CallFrescoBot/pkg/messages"
	"log"
	"os"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
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

			command, err := commands.GetCommand(update.Message.Text)
			if err != nil {
				log.Printf(err.Error())
			}

			response := command.RunCommand()

			message := tg.NewMessage(update.Message.Chat.ID, response)
			message.ReplyToMessageID = update.Message.MessageID

			bot.Send(message)
		}
	}
}
