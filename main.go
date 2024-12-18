package main

import (
	"CallFrescoBot/app"
	"CallFrescoBot/pkg/commands"
	"CallFrescoBot/pkg/consts"
	invoiceRepository "CallFrescoBot/pkg/repositories/invoice"
	callbackService "CallFrescoBot/pkg/service/callback"
	messageService "CallFrescoBot/pkg/service/message"
	"CallFrescoBot/pkg/service/numericKeyboard"
	planService "CallFrescoBot/pkg/service/plan"
	subsciptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	"bytes"
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
		if update.Message == nil && update.CallbackQuery == nil && update.PreCheckoutQuery == nil {
			continue
		}

		go func(upd tg.Update) {
			if err := handleUpdate(upd, bot); err != nil {
				log.Printf("Error handling update: %v", err)
			}
		}(update)
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

	messageErr := processMessage(update, bot, messageInfo)
	if messageErr != nil {
		sendMsgErr := messageService.SendMsgToUser(update.Message.Chat.ID, utils.LocalizeSafe(consts.ErrorMsg))
		if sendMsgErr != nil {
			return sendMsgErr
		}

		return fmt.Errorf("process message error: %w", messageErr)
	}

	callbackErr := processCallback(update, bot, messageInfo)
	if callbackErr != nil {
		return fmt.Errorf("process callback error: %w", callbackErr)
	}

	preCheckoutErr := processPreCheckout(update, bot)
	if preCheckoutErr != nil {
		return fmt.Errorf("process preCheckout error: %w", callbackErr)
	}

	successfulPaymentErr := handleSuccessfulPayment(update, bot)
	if successfulPaymentErr != nil {
		return fmt.Errorf("handle successful payment error: %w", successfulPaymentErr)
	}

	return nil
}

func handleSuccessfulPayment(update tg.Update, bot *tg.BotAPI) error {
	if update.Message == nil || update.Message.SuccessfulPayment == nil {
		return nil
	}

	db, _ := utils.GetDatabaseConnection()

	successfulPayment := update.Message.SuccessfulPayment
	log.Printf("Received SuccessfulPayment: %v", successfulPayment)

	payload := successfulPayment.InvoicePayload
	planID, err := extractPlanID(payload)
	if err != nil {
		log.Printf("Invalid payload in SuccessfulPayment: %v", err)
		return err
	}

	user, err := userService.GerUserByTgId(update.Message.From.ID)

	plan, err := planService.GetPlanById(planID)
	if err != nil {
		log.Printf("Failed to get plan by id: %v", err)
	}
	planId := plan.Id

	_, err = subsciptionService.CreateWithPlan(user, plan)

	_, _ = invoiceRepository.InvoiceCreate(
		1,
		user.Id,
		float64(update.Message.SuccessfulPayment.TotalAmount),
		update.Message.SuccessfulPayment.Currency,
		0,
		&planId,
		1,
		successfulPayment.TelegramPaymentChargeID,
		db,
	)

	confirmation := tg.NewMessage(update.Message.From.ID, utils.LocalizeSafe(consts.SubscriptionSuccess))
	_, err = bot.Send(confirmation)
	if err != nil {
		log.Printf("Failed to send confirmation message: %v", err)
		return err
	}

	return nil
}

func extractPlanID(payload string) (uint64, error) {
	var planID uint64
	n, err := fmt.Sscanf(payload, "plan_%d", &planID)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("invalid payload format")
	}
	return planID, nil
}

func processMessage(update tg.Update, bot *tg.BotAPI, messageInfo string) error {
	if update.Message == nil {
		return nil
	}

	if update.Message.SuccessfulPayment != nil {
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

	if update.Message.Voice != nil {
		transcription, err := handleVoiceMessage(update.Message, bot)
		if err != nil {
			return logAndNotifyOnErr(messageInfo, err)
		}
		update.Message.Text = transcription
	}

	responses, commandErr := commands.GetCommand(update, user).RunCommand()
	if notifyErr := logAndNotifyOnErr(messageInfo, commandErr); notifyErr != nil {
		return notifyErr
	}

	return sendBotResponses(bot, responses)
}

func handleVoiceMessage(message *tg.Message, bot *tg.BotAPI) (string, error) {
	fileURL, err := bot.GetFileDirectURL(message.Voice.FileID)
	if err != nil {
		return "", fmt.Errorf("failed to get file URL: %w", err)
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to download voice message: %w", err)
	}
	defer resp.Body.Close()

	oggPath := fmt.Sprintf("voice_%d.ogg", message.MessageID)
	oggFile, err := os.Create(oggPath)
	if err != nil {
		return "", fmt.Errorf("failed to create ogg file: %w", err)
	}
	defer oggFile.Close()

	_, err = io.Copy(oggFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save ogg file: %w", err)
	}

	wavPath := fmt.Sprintf("voice_%d.wav", message.MessageID)
	cmd := exec.Command("ffmpeg", "-i", oggPath, "-ar", "16000", "-ac", "1", wavPath)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	transcription, err := transcribeAudio(wavPath)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	os.Remove(oggPath)
	os.Remove(wavPath)

	return transcription, nil
}

func transcribeAudio(wavPath string) (string, error) {
	apiKey := utils.GetEnvVar("GPT_API_KEY")

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	file, err := os.Open(wavPath)
	if err != nil {
		return "", fmt.Errorf("failed to open wav file: %w", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(wavPath))
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	writer.WriteField("model", "whisper-1")

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/transcriptions", &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("transcription API error: %s", string(bodyBytes))
	}

	var result struct {
		Text string `json:"text"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to decode transcription response: %w", err)
	}

	return result.Text, nil
}

func sendBotResponses(bot *tg.BotAPI, responses []tg.Chattable) error {
	if responses == nil {
		return nil
	}
	for _, response := range responses {
		if _, err := bot.Send(response); err != nil {
			return err
		}
	}
	return nil
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

func processPreCheckout(update tg.Update, bot *tg.BotAPI) error {
	if update.PreCheckoutQuery == nil {
		return nil
	}

	preCheckout := update.PreCheckoutQuery

	log.Printf("Received PreCheckoutQuery with payload: %s", preCheckout.InvoicePayload)

	answer := tg.PreCheckoutConfig{
		PreCheckoutQueryID: preCheckout.ID,
		OK:                 true,
		ErrorMessage:       "",
	}

	_, err := bot.Request(answer)
	if err != nil {
		log.Printf("Failed to answer PreCheckoutQuery: %v", err)
		return fmt.Errorf("answer pre_checkout_query error: %w", err)
	}

	log.Println("PreCheckoutQuery answered successfully with ok=true")
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
	}
	return nil
}
