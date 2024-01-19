package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"CallFrescoBot/pkg/utils"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type StatusCommand struct {
	BaseCommand
}

func (cmd StatusCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return tg.NewMessage(cmd.Update.Message.Chat.ID, result), err
	}

	subscription, err := subscriptionService.GetUserSubscriptionWithNoPlanLimit(cmd.User)
	if err != nil {
		return nil, err
	}

	subscriptionName := ResolveSubscriptionName(subscription.Limit)
	messagesCount, err := messageService.CountMessagesByUserAndDate(cmd.User, subscription.Limit, time.Now().AddDate(0, 0, -1))
	if err != nil {
		return nil, err
	}

	remainingMessages := RemainingMessages(int64(subscription.Limit), messagesCount)
	validDue := SubscriptionValidDue(subscription)

	status := fmt.Sprintf(utils.LocalizeSafe(consts.StatusMsg), subscriptionName, subscription.Limit, remainingMessages, validDue)

	return tg.NewMessage(cmd.Update.Message.Chat.ID, status), nil
}

func RemainingMessages(subscriptionLimit int64, messagesCount int64) int64 {
	result := subscriptionLimit - messagesCount
	if result < 0 {
		return 0
	}

	return subscriptionLimit - messagesCount
}

func ResolveSubscriptionName(limit int) string {
	switch limit := limit; {
	case limit == 0:
		return utils.LocalizeSafe(consts.SubscriptionPlanHacker)
	case limit <= 5:
		return utils.LocalizeSafe(consts.SubscriptionPlanFree)
	case limit <= 10:
		return utils.LocalizeSafe(consts.SubscriptionPlanStart)
	case limit <= 20:
		return utils.LocalizeSafe(consts.SubscriptionPlanPro)
	case limit <= 50:
		return utils.LocalizeSafe(consts.SubscriptionPlanBoss)
	default:
		return utils.LocalizeSafe(consts.SubscriptionPlanHacker)
	}
}

func SubscriptionValidDue(subscription *models.Subscription) string {
	if subscription.ActiveDue.IsZero() {
		return "❌ No active subscriptions"
	}

	return subscription.ActiveDue.Format("02.01.2006")
}
