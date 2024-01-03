package main

import (
	"CallFrescoBot/app"
	"CallFrescoBot/pkg/commands"
	messageService "CallFrescoBot/pkg/service/message"
	UserService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	app.SetupApp()

	upd := tg.NewUpdate(0)
	upd.Timeout = 15

	bot := utils.GetBot()
	updates := bot.GetUpdatesChan(upd)

	for update := range updates {
		if update.Message != nil {
			messageInfo := fmt.Sprintf("[%s, %d] %s",
				update.Message.From.UserName,
				update.Message.Chat.ID,
				update.Message.Text,
			)
			log.Printf(messageInfo)
			user, err := UserService.GetOrCreate(update)
			if err != nil {
				log.Printf(err.Error())
				err := messageService.SendErrorToUser(messageInfo, err)
				if err != nil {
					log.Printf(err.Error())
				}
			}

			response, err := commands.GetCommand(update, user).RunCommand()
			if err != nil {
				log.Printf(err.Error())
				err := messageService.SendErrorToUser(messageInfo, err)
				if err != nil {
					log.Printf(err.Error())
				}
			}

			if response != nil {
				_, err = bot.Send(response)
				if err != nil {
					log.Printf(err.Error())
				}
			}
		}
	}
}
