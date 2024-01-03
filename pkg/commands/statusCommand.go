package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type StatusCommand struct {
	Update tg.Update
	User   *models.User
}

func (cmd StatusCommand) Common() (string, error) {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Update.Message.Text)
	if err != nil {
		return messageValidatorText, err
	}

	return "", nil
}

func (cmd StatusCommand) RunCommand() (tg.Chattable, error) {
	result, err := cmd.Common()
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

	status := fmt.Sprintf(consts.StatusMsg, subscriptionName, subscription.Limit, remainingMessages, validDue)

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
		return consts.SubscriptionPlanHacker
	case limit <= 5:
		return consts.SubscriptionPlanHomeless
	case limit <= 25:
		return consts.SubscriptionPlanBasic
	case limit <= 50:
		return consts.SubscriptionPlanVIP
	case limit <= 100:
		return consts.SubscriptionPlanDeluxe
	default:
		return consts.SubscriptionPlanHacker
	}
}

func SubscriptionValidDue(subscription *models.Subscription) string {
	if subscription.ActiveDue.IsZero() {
		return "❌ No active subscriptions"
	}

	return subscription.ActiveDue.Format("02.01.2006")
}
