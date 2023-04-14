package utils

import (
	"CallFrescoBot/pkg/consts"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

var apiKey string
var botDebug string
var botApi *tg.BotAPI

func init() {
	apiKey = GetEnvVar("TELEGRAM_API_KEY")
	botDebug = GetEnvVar("BOT_DEBUG")
}

func CreateBot() error {
	if apiKey == "" {
		log.Fatalln(consts.MissingTgKey)
	}

	bot, err := tg.NewBotAPI(apiKey)
	if err != nil {
		log.Panic(err)
	}

	bDebug, _ := strconv.ParseBool(botDebug)

	bot.Debug = bDebug
	botApi = bot

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return err
}

func GetBot() *tg.BotAPI {
	return botApi
}
