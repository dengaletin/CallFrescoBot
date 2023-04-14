package messageService

import (
	"CallFrescoBot/pkg/consts"
	messageRepository "CallFrescoBot/pkg/repositories/message"
	"CallFrescoBot/pkg/utils"
	"errors"
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
