package userService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	subscriptionRepository "CallFrescoBot/pkg/repositories/subscription"
	"CallFrescoBot/pkg/repositories/user"
	"CallFrescoBot/pkg/utils"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func GetOrCreate(tgUser *tg.User) (*models.User, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	user, err := userRepository.FirstOrCreate(tgUser, db)
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

	subscriptionLimit := subscriptionRepository.GetUserSubscriptionLimit(user, db)
	messagesCount, err := messageRepository.CountMessagesByUserAndDate(user, subscriptionLimit, time.Now().AddDate(0, 0, -1), db)
	if err != nil {
		return "", err
	}

	if int64(subscriptionLimit) == 0 {
		return "", nil
	}

	if messagesCount >= int64(subscriptionLimit) {
		return consts.RunOutOfMessages, errors.New("out of messages")
	}

	return "", nil
}
