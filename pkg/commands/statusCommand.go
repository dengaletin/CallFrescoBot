package commands

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageService "CallFrescoBot/pkg/service/message"
	planService "CallFrescoBot/pkg/service/plan"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"CallFrescoBot/pkg/types"
	"CallFrescoBot/pkg/utils"
	"encoding/json"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

type StatusCommand struct {
	BaseCommand
}

func (cmd StatusCommand) RunCommand() ([]tg.Chattable, error) {
	result, err := cmd.Common(false)
	if err != nil {
		return []tg.Chattable{tg.NewMessage(cmd.Update.Message.Chat.ID, result)}, err
	}

	subscription, err := subscriptionService.GetUserSubscriptionWithNoPlanLimit(cmd.User)
	if err != nil {
		return nil, err
	}

	var msgText string

	if subscription == nil {
		msgText = utils.LocalizeSafe(consts.FreeSubscriptionFinish)
	} else if subscription.PlanId != nil {
		plan, err := planService.GetPlanById(*subscription.PlanId)
		if err != nil {
			return nil, err
		}

		var usage types.Usage

		err = json.Unmarshal(subscription.Usage, &usage)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling usage JSON: %w", err)
		}

		var config types.Config
		if err := json.Unmarshal(plan.Config, &config); err != nil {
			return nil, err
		}

		limitsInfo := ""
		if config.Limit.Gpt4OMiniLimit > 0 {
			limitsInfo += fmt.Sprintf(
				"*GPT4o-mini* - %d / %d "+utils.LocalizeSafe(consts.Requests)+"\n",
				usage.Gpt4OMini+usage.Gpt4OMiniContext, config.Limit.Gpt4OMiniLimit,
			)
		}
		if config.Limit.Gpt4OLimit > 0 {
			limitsInfo += fmt.Sprintf(
				"*GPT4o* - %d / %d "+utils.LocalizeSafe(consts.Requests)+"\n",
				usage.Gpt4O+usage.Gpt4OContext, config.Limit.Gpt4OLimit,
			)
		}
		if config.Limit.Dalle3Limit > 0 {
			limitsInfo += fmt.Sprintf(
				"*Dalle3* - %d / %d "+utils.LocalizeSafe(consts.Requests)+"\n",
				usage.Dalle3+usage.Dalle3Context, config.Limit.Dalle3Limit,
			)
		}
		if config.Limit.Gpt4O1Limit > 0 {
			limitsInfo += fmt.Sprintf(
				"*GPT4o1* - %d / %d "+utils.LocalizeSafe(consts.Requests)+"\n",
				usage.Gpt4O1+usage.Gpt4O1Context, config.Limit.Gpt4O1Limit,
			)
		}
		contextSupport := utils.LocalizeSafe(consts.No)
		if config.Limit.ContextSupport {
			contextSupport = utils.LocalizeSafe(consts.Yes)
		}
		limitsInfo += fmt.Sprintf("*"+utils.LocalizeSafe(consts.ContextSupport)+"* - %s\n", contextSupport)

		validDue := SubscriptionValidDue(subscription)
		msgText = fmt.Sprintf(utils.LocalizeSafe(consts.PlanStatusMsg), plan.Name, limitsInfo, validDue)
	} else {
		subscriptionName := ResolveSubscriptionName(subscription.Limit)
		messagesCount, err := messageService.CountMessagesByUserAndDate(cmd.User, subscription.Limit, time.Now().AddDate(0, 0, -1))
		if err != nil {
			return nil, err
		}

		remainingMessages := RemainingMessages(int64(subscription.Limit), messagesCount)
		validDue := SubscriptionValidDue(subscription)
		msgText = fmt.Sprintf(utils.LocalizeSafe(consts.StatusMsg), subscriptionName, subscription.Limit, remainingMessages, validDue)
	}

	msg := tg.NewMessage(cmd.Update.Message.Chat.ID, msgText)
	msg.ParseMode = "markdown"

	return []tg.Chattable{msg}, nil
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
	case limit <= consts.NoPlanLimit:
		return utils.LocalizeSafe(consts.SubscriptionPlanFree)
	default:
		return utils.LocalizeSafe(consts.SubscriptionPlanHacker)
	}
}

func SubscriptionValidDue(subscription *models.Subscription) string {
	return subscription.ActiveDue.Format("02.01.2006")
}
