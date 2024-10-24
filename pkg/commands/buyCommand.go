package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/utils"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type BuyCommand struct {
	BaseCommand
}

func (cmd BuyCommand) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	startStarInvoice(cmd)

	return []tg.Chattable{}, nil
}

func startStarInvoice(cmd BuyCommand) {
	price := []tg.LabeledPrice{
		{
			Label:  "XTR",
			Amount: 300,
		},
	}

	invoice := tg.NewInvoice(
		cmd.Update.Message.Chat.ID,
		"Fresco AI Premium",
		utils.LocalizeSafe(consts.BuyMsg),
		"plan_19",
		"",
		"start_fresco_ai_unique",
		"XTR",
		price,
	)
	invoice.SuggestedTipAmounts = []int{}

	bot := utils.GetBot()
	sentInvoice, err := bot.Send(invoice)
	if err != nil {
		log.Printf("Error sending invoice: %v", err)
		return
	}
	fmt.Println("Invoice sent successfully:", sentInvoice)
}
