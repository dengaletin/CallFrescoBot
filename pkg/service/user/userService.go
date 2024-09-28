package userService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	userRepository "CallFrescoBot/pkg/repositories/user"
	planService "CallFrescoBot/pkg/service/plan"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"CallFrescoBot/pkg/types"
	"CallFrescoBot/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

func dbConnection() (*gorm.DB, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}
	return db, nil
}

func GetOrCreate(user *tgbotapi.User) (*models.User, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}

	botUser, err := userRepository.FirstOrCreate(user, db)
	if err != nil {
		return nil, err
	}

	if botUser.IsNew {
		_, err := subscriptionService.Create(botUser, consts.NoPlanLimit)
		if err != nil {
			return nil, err
		}
	}

	return botUser, nil
}

func ValidateUser(user *models.User) (string, error) {
	db, err := dbConnection()
	if err != nil {
		return "", err
	}
	if !user.IsActive {
		return "sorry, your profile is not active", errors.New("profile is not active")
	}
	subscription, err := subscriptionService.GetUserSubscriptionWithNoPlanLimit(user)
	if err != nil {
		return "", err
	}

	if subscription == nil {
		return utils.LocalizeSafe(consts.FreeSubscriptionFinish), errors.New("free subscription finish")
	} else if subscription.PlanId != nil {
		var usage types.Usage
		var usageCount int
		var currentLimit int
		var config types.Config

		err = json.Unmarshal(subscription.Usage, &usage)
		if err != nil {
			return "", fmt.Errorf("error unmarshaling usage JSON: %w", err)
		}

		userPlan, err := planService.GetPlanById(*subscription.PlanId)
		if err != nil {
			return "", fmt.Errorf("can't get subscription plan: %w", err)
		}

		if err := json.Unmarshal(userPlan.Config, &config); err != nil {
			return "", err
		}

		userMode := user.Mode

		if user.Dialog != 0 && !config.Limit.ContextSupport {
			return "dialog context is not supported in your subscription", errors.New("dialog context is not supported in your subscription")
		}

		switch userMode {
		case consts.UsageModeGpt4OMini:
			if config.Limit.Gpt4OMiniLimit <= 0 {
				return "Model is not available in your subscription", errors.New("GPT-4o-mini model is not available in your subscription")
			}
			usageCount = usage.Gpt4OMini + usage.Gpt4OMiniContext
			currentLimit = config.Limit.Gpt4OMiniLimit
		case consts.UsageModeDalle3:
			if config.Limit.Dalle3Limit <= 0 {
				return "Model is not available in your subscription", errors.New("DALL-E-3 model is not available in your subscription")
			}
			usageCount = usage.Dalle3 + usage.Dalle3Context
			currentLimit = config.Limit.Dalle3Limit
		case consts.UsageModeGpt4O:
			if config.Limit.Gpt4OLimit <= 0 {
				return "Model is not available in your subscription", errors.New("GPT-4-o model is not available in your subscription")
			}
			usageCount = usage.Gpt4O + usage.Gpt4OContext
			currentLimit = config.Limit.Gpt4OLimit
		case consts.UsageModeGpt4O1:
			if config.Limit.Gpt4O1Limit <= 0 {
				return "Model is not available in your subscription", errors.New("GPT-4-o1 model is not available in your subscription")
			}
			usageCount = usage.Gpt4O1 + usage.Gpt4O1Context
			currentLimit = config.Limit.Gpt4O1Limit
		default:
			return "", fmt.Errorf("unknown usage mode: %w", err)
		}

		if usageCount >= currentLimit {
			return utils.LocalizeSafe(consts.RunOutOfMessages), errors.New("out of messages")
		}
	} else {
		limitDate := time.Now().AddDate(0, 0, -1)
		messagesCount, err := messageRepository.CountMessagesByUserAndDate(user, subscription.Limit, limitDate, db)
		if err != nil {
			return "", err
		}

		if messagesCount >= int64(subscription.Limit) {
			return utils.LocalizeSafe(consts.RunOutOfMessages), errors.New("out of messages")
		}
	}
	return "", nil
}

func GerUserByTgId(tgId int64) (*models.User, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}
	return userRepository.GerUserByTgId(tgId, db)
}

func GetUserById(id uint64) (*models.User, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}
	return userRepository.GetUserById(id, db)
}

func SetMode(mode int64, user *models.User) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	return userRepository.SetMode(mode, user, db)
}

func SetLanguage(language int64, user *models.User) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}

	return userRepository.SetLanguage(language, user, db)
}

func SetDialogStatus(dialogStatus int64, user *models.User) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}
	return userRepository.SetDialogStatus(dialogStatus, user, db)
}

func GetMode(mode int64) (string, error) {
	return userRepository.GetMode(mode)
}

func SetUserDialogFromId(user *models.User) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}

	message, err := messageRepository.GetUserLastMessage(user, db)
	if err != nil {
		return err
	}

	return userRepository.SetDialogFromId(message.Id, user, db)
}
