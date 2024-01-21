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
	"net/http"
)

func main() {
	app.SetupApp()
	go setupServer()
	updates := initBotUpdates(15)
	processUpdates(updates)
}

func setupServer() {
	http.HandleFunc("/payment-callback", postHandler)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintln(w, "OK")
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
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

		if err := handleUpdate(update, bot); err != nil {
			log.Printf("Error handling update: %v", err)
		}
	}
}

func handleUpdate(update tg.Update, bot *tg.BotAPI) error {
	var messageInfo string
	if update.Message != nil {
		messageInfo = formatMessageInfo(update.Message)
		log.Printf(messageInfo)
	}

	if update.CallbackQuery != nil {
		messageInfo = formatMessageInfo(update.CallbackQuery.Message)
		log.Printf(messageInfo)
	}

	if err := processMessage(update, bot, messageInfo); err != nil {
		return fmt.Errorf("process message error: %w", err)
	}

	if err := processCallback(update, bot, messageInfo); err != nil {
		return fmt.Errorf("process callback error: %w", err)
	}

	return nil
}

func processMessage(update tg.Update, bot *tg.BotAPI, messageInfo string) error {
	if update.Message == nil {
		return nil
	}

	_, from, messageServiceErr := messageService.ParseUpdate(update)
	if err := logAndNotifyOnErr("", messageServiceErr); err != nil {
		return err
	}

	user, userServiceErr := userService.GetOrCreate(from)
	if err := logAndNotifyOnErr(messageInfo, userServiceErr); err != nil {
		return err
	}

	utils.InitBundle(user.Lang)

	if mainMenuErr := numericKeyboard.CreateMainMenu(); mainMenuErr != nil {
		return logAndNotifyOnErr(messageInfo, mainMenuErr)
	}

	response, commandErr := commands.GetCommand(update, user).RunCommand()
	if err := logAndNotifyOnErr(messageInfo, commandErr); err != nil {
		return err
	}

	return sendBotResponse(bot, response)
}

func processCallback(update tg.Update, bot *tg.BotAPI, messageInfo string) error {
	if update.CallbackQuery == nil {
		return nil
	}

	fmt.Println(update.CallbackQuery.Data)
	user, userServiceErr := userService.GetOrCreate(update.CallbackQuery.From)
	if userServiceErr != nil {
		return fmt.Errorf("get user error: %w", userServiceErr)
	}

	utils.InitBundle(user.Lang)

	callbackErr := callbackService.ResolveAndHandle(update.CallbackQuery, user, bot)
	if err := logAndNotifyOnErr(messageInfo, callbackErr); err != nil {
		return err
	}

	return nil
}

func formatMessageInfo(message *tg.Message) string {
	return fmt.Sprintf(
		"[%s, %d] %s",
		message.From.UserName,
		message.Chat.ID,
		message.Text,
	)
}

func logAndNotifyOnErr(messageInfo string, err error) error {
	if err != nil {
		log.Printf(err.Error())
		errMsg := fmt.Sprintf("❌❌❌ Error: [%s] %s", messageInfo, err.Error())
		if notifyErr := messageService.SendMsgToUser(consts.LogErrorRecipient, errMsg); notifyErr != nil {
			log.Printf(notifyErr.Error())
			return fmt.Errorf("error sending notification: %w", notifyErr)
		}
		return err
	}
	return nil
}

func sendBotResponse(bot *tg.BotAPI, response tg.Chattable) error {
	if response == nil {
		return nil
	}
	_, err := bot.Send(response)

	return err
}
