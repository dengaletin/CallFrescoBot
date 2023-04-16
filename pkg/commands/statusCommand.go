package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"fmt"
	"log"
	"time"
)

type StatusCommand struct {
	Message string
	User    *models.User
}

func (cmd StatusCommand) Common() string {
	messageValidatorText, err := messageService.ValidateMessage(cmd.Message)
	if err != nil {
		log.Printf(err.Error())
		return messageValidatorText
	}

	return ""
}

func (cmd StatusCommand) RunCommand() string {
	result := cmd.Common()

	if result != "" {
		return result
	}

	subscription, err := subscriptionService.GetUserSubscriptionWithNoPlanLimit(cmd.User)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}

	subscriptionName := ResolveSubscriptionName(subscription.Limit)
	messagesCount, err := messageService.CountMessagesByUserAndDate(cmd.User, subscription.Limit, time.Now().AddDate(0, 0, -1))
	if err != nil {
		log.Printf(err.Error())
		return ""
	}

	remainingMessages := RemainingMessages(int64(subscription.Limit), messagesCount)
	validDue := SubscriptionValidDue(subscription)

	status := fmt.Sprintf(consts.StatusMsg, subscriptionName, subscription.Limit, remainingMessages, validDue)

	return status
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
	case limit <= 15:
		return consts.SubscriptionPlanBomj
	case limit <= 50:
		return consts.SubscriptionPlanStudent
	case limit <= 100:
		return consts.SubscriptionPlanMajor
	case limit <= 200:
		return consts.SubscriptionPlanGigaSheikh
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
