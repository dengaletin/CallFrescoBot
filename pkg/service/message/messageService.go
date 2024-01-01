package messageService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	"CallFrescoBot/pkg/utils"
	"errors"
	"time"
)

func ValidateMessage(cmd string) (string, error) {
	if cmd == "" {
		return consts.UnsupportedMessageType, errors.New("unsupported message type")
	}
	if len([]rune(cmd)) < 4 {
		return consts.MessageIsTooShort, errors.New("message is too short")
	}

	return cmd, nil
}

func CreateMessage(userId uint64, message string, response string) error {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return errors.New("error occurred while getting a DB connection from the connection pool")
	}

	_, err = messageRepository.MessageCreate(userId, message, response, db)
	if err != nil {
		return err
	}

	return nil
}

func CountMessagesByUserAndDate(user *models.User, limit int, date time.Time) (int64, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return 0, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	messagesCount, err := messageRepository.CountMessagesByUserAndDate(user, limit, date, db)

	if err != nil {
		return 0, err
	}

	return messagesCount, nil
}

func GetMessagesByUser(user *models.User, limit int) ([]models.Message, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	messages, err := messageRepository.GetMessagesByUser(user, limit, db)

	if err != nil {
		return nil, err
	}

	return messages, nil
}
