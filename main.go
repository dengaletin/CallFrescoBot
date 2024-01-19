package main

import (
	"CallFrescoBot/app"
	"CallFrescoBot/pkg/commands"
	"CallFrescoBot/pkg/consts"
	callbackService "CallFrescoBot/pkg/service/callback"
	messageService "CallFrescoBot/pkg/service/message"
	"CallFrescoBot/pkg/service/numericKeyboard"
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
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		message, from, messageServiceErr := messageService.ParseUpdate(update)
		if messageServiceErr != nil {
			logAndNotify("", messageServiceErr)
		}

		messageInfo := formatMessageInfo(message)
		log.Printf(messageInfo)

		user, userServiceErr := userService.GetOrCreate(from)
		if userServiceErr != nil {
			logAndNotify(messageInfo, userServiceErr)
		}

		locErr := utils.CreateLoc(user)
		if locErr != nil {
			log.Printf(locErr.Error())
			logAndNotify(messageInfo, locErr)
		}

		mainMenuErr := numericKeyboard.CreateMainMenu()
		if mainMenuErr != nil {
			log.Printf(mainMenuErr.Error())
			logAndNotify(messageInfo, mainMenuErr)
		}

		if update.CallbackQuery != nil {
			fmt.Println(update.CallbackQuery.Data)
			callbackErr := callbackService.ResolveAndHandle(update.CallbackQuery, user, bot)
			if callbackErr != nil {
				log.Printf(callbackErr.Error())
				logAndNotify(messageInfo, callbackErr)
			}
		}

		if update.Message == nil {
			continue
		}

		response, commandErr := commands.GetCommand(update, user).RunCommand()
		if commandErr != nil {
			logAndNotify(messageInfo, commandErr)
		}

		if err := sendBotResponse(bot, response); err != nil {
			log.Printf(err.Error())
			logAndNotify(messageInfo, err)
		}
	}
}

func formatMessageInfo(message *tg.Message) string {
	return fmt.Sprintf(
		"[%s, %d] %s",
		message.From.UserName,
		message.Chat.ID,
		message.Text,
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
