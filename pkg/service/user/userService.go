package userService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	userRepository "CallFrescoBot/pkg/repositories/user"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"CallFrescoBot/pkg/utils"
	"errors"
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

	return userRepository.FirstOrCreate(user, db)
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

	limitDate := time.Now().AddDate(0, 0, -1)
	messagesCount, err := messageRepository.CountMessagesByUserAndDate(user, subscription.Limit, limitDate, db)
	if err != nil {
		return "", err
	}

	if messagesCount >= int64(subscription.Limit) {
		return utils.LocalizeSafe(consts.RunOutOfMessages), errors.New("out of messages")
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
