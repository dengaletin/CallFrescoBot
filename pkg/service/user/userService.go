package userService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	"CallFrescoBot/pkg/repositories/user"
	subscriptionService "CallFrescoBot/pkg/service/subsciption"
	"CallFrescoBot/pkg/utils"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func GetOrCreate(cmd tgbotapi.Update) (*models.User, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	user, err := userRepository.FirstOrCreate(cmd.Message.From, db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func ValidateUser(user *models.User) (string, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return "", errors.New("error occurred while getting a DB connection from the connection pool")
	}

	if user.IsActive != true {
		return "Sorry man, your profile is not active", errors.New("profile is not active")
	}

	subscription, err := subscriptionService.GetUserSubscriptionWithNoPlanLimit(user)
	if err != nil {
		return "", err
	}

	messagesCount, err := messageRepository.CountMessagesByUserAndDate(user, subscription.Limit, time.Now().AddDate(0, 0, -1), db)
	if err != nil {
		return "", err
	}

	if int64(subscription.Limit) == 0 {
		return "", nil
	}

	if messagesCount >= int64(subscription.Limit) {
		return consts.RunOutOfMessages, errors.New("out of messages")
	}

	return "", nil
}

func GerUserByTgId(tgId int64) (*models.User, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	user, err := userRepository.GerUserByTgId(tgId, db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func SetMode(mode int64, user *models.User) error {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return errors.New("error occurred while getting a DB connection from the connection pool")
	}

	err = userRepository.SetMode(mode, user, db)
	if err != nil {
		return err
	}

	return nil
}

func SetDialogStatus(dialogStatus int64, user *models.User) error {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return errors.New("error occurred while getting a DB connection from the connection pool")
	}

	err = userRepository.SetDialogStatus(dialogStatus, user, db)
	if err != nil {
		return err
	}

	return nil
}

func GetMode(mode int64) (string, error) {
	result, err := userRepository.GetMode(mode)
	if err != nil {
		return "", err
	}

	return result, nil
}

func GetDialogStatus(dialogStatus int64) (string, error) {
	result, err := userRepository.GetDialogStatus(dialogStatus)
	if err != nil {
		return "", err
	}

	return result, nil
}
