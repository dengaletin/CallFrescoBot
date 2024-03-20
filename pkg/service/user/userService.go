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
		_, err := subscriptionService.Create(botUser, 3)
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
		case consts.UsageModeGpt35:
			if config.Limit.Gpt35Limit <= 0 {
				return "GPT-3.5 model is not available in your subscription", errors.New("GPT-3.5 model is not available in your subscription")
			}
			usageCount = usage.Gpt35 + usage.Gpt35Context
			currentLimit = config.Limit.Gpt35Limit
		case consts.UsageModeDalle3:
			if config.Limit.Dalle3Limit <= 0 {
				return "DALL-E 3 model is not available in your subscription", errors.New("DALL-E 3 model is not available in your subscription")
			}
			usageCount = usage.Dalle3 + usage.Dalle3Context
			currentLimit = config.Limit.Dalle3Limit
		case consts.UsageModeGpt4:
			if config.Limit.Gpt4Limit <= 0 {
				return "GPT-4 model is not available in your subscription", errors.New("GPT-4 model is not available in your subscription")
			}
			usageCount = usage.Gpt4 + usage.Gpt4Context
			currentLimit = config.Limit.Gpt4Limit
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
