package messageService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	"CallFrescoBot/pkg/utils"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

func ValidateMessage(cmd string) (string, error) {
	if cmd == "" {
		return consts.UnsupportedMessageType, errors.New(consts.UnsupportedMessageType)
	}
	if len([]rune(cmd)) < 4 {
		return consts.MessageIsTooShort, errors.New(consts.MessageIsTooShort)
	}

	return cmd, nil
}

func dbConnection() (*gorm.DB, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("failed to obtain a db connection from the pool")
	}

	return db, nil
}

func CreateMessage(userId uint64, message string, response string) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}

	_, err = messageRepository.MessageCreate(userId, message, response, db)
	return err
}

func CountMessagesByUserAndDate(user *models.User, limit int, date time.Time) (int64, error) {
	db, err := dbConnection()
	if err != nil {
		return 0, err
	}

	return messageRepository.CountMessagesByUserAndDate(user, limit, date, db)
}

func GetMessagesByUser(user *models.User, limit int) ([]models.Message, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}

	return messageRepository.GetMessagesByUser(user, limit, db)
}

func SendMsgToUser(chatId int64, msgInfo string) error {
	message := tg.NewMessage(chatId, msgInfo)
	bot := utils.GetBot()

	_, err := bot.Send(message)
	if err != nil {
		return err
	}

	return nil
}
