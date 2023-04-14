package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	subscriptionRepository "CallFrescoBot/pkg/repositories/subscription"
	"CallFrescoBot/pkg/utils"
	"fmt"
	"log"
	"time"
)

type StatusCommand struct {
	Message string
	User    *models.User
}

func (cmd StatusCommand) RunCommand() string {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		log.Printf(err.Error())
		return ""
	}

	subscriptionLimit := subscriptionRepository.GetUserSubscriptionLimit(cmd.User, db)
	subscriptionName := ResolveSubscriptionName(subscriptionLimit)
	messagesCount, err := messageRepository.CountMessagesByUserAndDate(cmd.User, subscriptionLimit, time.Now().AddDate(0, 0, -1), db)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}

	remainingMessages := RemainingMessages(int64(subscriptionLimit), messagesCount)

	status := fmt.Sprintf(consts.StatusMsg, subscriptionName, subscriptionLimit, remainingMessages)

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
