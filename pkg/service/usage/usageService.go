package usageService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"CallFrescoBot/pkg/types"
	"encoding/json"
	"fmt"
)

func SaveUsage(user *models.User, modeId int64) error {
	subscription, err := subscriptionService.GetUserSubscription(user)
	if err != nil {
		return fmt.Errorf("error getting user subscription: %w", err)
	}
	var usage types.Usage

	err = json.Unmarshal(subscription.Usage, &usage)
	if err != nil {
		return fmt.Errorf("error unmarshaling usage JSON: %w", err)
	}

	switch modeId {
	case consts.UsageModeGpt4OMini:
		usage.Gpt4OMini += 1
	case consts.UsageModeDalle3:
		usage.Dalle3 += 1
	case consts.UsageModeGpt4O:
		usage.Gpt4O += 1
	case consts.UsageModeGpt4O1:
		usage.Gpt4O1 += 1
	case consts.UsageModeGpt4OMiniContext:
		usage.Gpt4OMiniContext += 1
	case consts.UsageModeDalle3Context:
		usage.Dalle3Context += 1
	case consts.UsageModeGpt4OContext:
		usage.Gpt4OContext += 1
	case consts.UsageModeGpt4O1Context:
		usage.Gpt4O1Context += 1
	case consts.UsageModeGpt4OMiniVoice:
		usage.Gpt4OMiniVoice += 1
	case consts.UsageModeDalle3Voice:
		usage.Dalle3Voice += 1
	case consts.UsageModeGpt4OVoice:
		usage.Gpt4OVoice += 1
	case consts.UsageModeGpt4O1Voice:
		usage.Gpt4O1Voice += 1
	case consts.UsageModeGpt4OMiniContextVoice:
		usage.Gpt4OMiniContextVoice += 1
	case consts.UsageModeDalle3ContextVoice:
		usage.Dalle3ContextVoice += 1
	case consts.UsageModeGpt4OContextVoice:
		usage.Gpt4OContextVoice += 1
	case consts.UsageModeGpt4O1ContextVoice:
		usage.Gpt4O1ContextVoice += 1
	default:
		return fmt.Errorf("unknown usage mode: %w", err)
	}

	updatedUsage, err := json.Marshal(usage)
	if err != nil {
		return fmt.Errorf("error marshaling updated usage to JSON: %w", err)
	}

	subscription.Usage = updatedUsage

	err = subscriptionService.UpdateUserSubscription(subscription)
	if err != nil {
		return fmt.Errorf("error updating user subscription: %w", err)
	}

	return nil
}
