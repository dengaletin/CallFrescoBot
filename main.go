package main

import (
	"CallFrescoBot/app"
	"CallFrescoBot/pkg/commands"
	"CallFrescoBot/pkg/consts"
	messageService "CallFrescoBot/pkg/service/message"
	userService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func main() {
	app.SetupApp()
	updates := initBotUpdates(15)
	processUpdates(updates)
}

func initBotUpdates(timeout int) tg.UpdatesChannel {
	upd := tg.NewUpdate(0)
	upd.Timeout = timeout

	bot := utils.GetBot()
	return bot.GetUpdatesChan(upd)
}

func processUpdates(updates tg.UpdatesChannel) {
	bot := utils.GetBot()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		messageInfo := formatMessageInfo(update)
		log.Printf(messageInfo)

		user, userServiceErr := userService.GetOrCreate(update)
		if userServiceErr != nil {
			logAndNotify(messageInfo, userServiceErr)
		}

		response, commandErr := commands.GetCommand(update, user).RunCommand()
		if commandErr != nil {
			logAndNotify(messageInfo, commandErr)
		}

		if err := sendBotResponse(bot, response); err != nil {
			log.Printf(err.Error())
		}
	}
}

func formatMessageInfo(update tg.Update) string {
	return fmt.Sprintf(
		"[%s, %d] %s",
		update.Message.From.UserName,
		update.Message.Chat.ID,
		update.Message.Text,
	)
}

func logAndNotify(messageInfo string, err error) {
	log.Printf(err.Error())
	errMsg := fmt.Sprintf("❌❌❌ Error: [%s] %s", messageInfo, err.Error())
	if notifyErr := messageService.SendMsgToUser(consts.LogErrorRecipient, errMsg); notifyErr != nil {
		log.Printf(notifyErr.Error())
	}
}

func sendBotResponse(bot *tg.BotAPI, response tg.Chattable) error {
	if response == nil {
		return nil
	}
	_, err := bot.Send(response)
	return err
}
