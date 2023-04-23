package commands

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type ICommand interface {
	RunCommand() (tg.Chattable, error)
	Common() (string, error)
}
