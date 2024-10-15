package messageService

import (
	"CallFrescoBot/pkg/models"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	"CallFrescoBot/pkg/utils"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

func ParseUpdate(update tg.Update) (*tg.Message, *tg.User, error) {
	if update.Message != nil {
		return update.Message, update.Message.From, nil
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message, update.CallbackQuery.From, nil
	}

	return nil, nil, errors.New("can't parse message")
}

func ValidateMessage(cmd string) (string, error) {
	// todo: remove?
	//if cmd == "" {
	//	return utils.LocalizeSafe(consts.UnsupportedMessageType), errors.New("unsupported message type")
	//}
	//if len([]rune(cmd)) < 4 {
	//	return utils.LocalizeSafe(consts.MessageIsTooShort), errors.New("message is too short")
	//}

	return cmd, nil
}

func dbConnection() (*gorm.DB, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("failed to obtain a db connection from the pool")
	}

	return db, nil
}

func CreateMessage(userId uint64, message string, response string, mode int64) error {
	db, err := dbConnection()
	if err != nil {
		return err
	}

	_, err = messageRepository.MessageCreate(userId, message, response, mode, db)
	return err
}

func CountMessagesByUserAndDate(user *models.User, limit int, date time.Time) (int64, error) {
	db, err := dbConnection()
	if err != nil {
		return 0, err
	}

	return messageRepository.CountMessagesByUserAndDate(user, limit, date, db)
}

func GetMessagesByUser(user *models.User, limit int, mode int64) ([]models.Message, error) {
	db, err := dbConnection()
	if err != nil {
		return nil, err
	}

	return messageRepository.GetMessagesByUser(user, limit, mode, db)
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
