package messageRepository

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func MessageCreate(userId uint64, message string, response string, mode int64, db *gorm.DB) (*models.Message, error) {
	newMessage := models.Message{
		UserId:   userId,
		Message:  message,
		Response: response,
		Mode:     mode,
	}
	result := db.Create(&newMessage)

	if result.Error != nil && result.RowsAffected != 1 {
		return nil, errors.New("error occurred while creating a new message")
	}

	return &newMessage, nil
}

func CountMessagesByUserAndDate(user *models.User, limit int, date time.Time, db *gorm.DB) (int64, error) {
	var messages []models.Message

	result := db.Where("user_id = ? AND created_at > ?", user.Id, date.String()).Find(&messages).Limit(limit)

	if result.Error != nil && result.RowsAffected != 1 {
		return 0, errors.New("error occurred while counting messages")
	}

	return result.RowsAffected, nil
}

func GetMessagesByUser(user *models.User, limit int, mode int64, db *gorm.DB) ([]models.Message, error) {
	var messages []models.Message
	var rawSQL string

	baseSQL := "SELECT * FROM messages WHERE user_id = %d AND id > %d AND mode = %d"

	if user.Id == consts.AdminUserId {
		rawSQL = fmt.Sprintf(baseSQL+" ORDER BY created_at DESC LIMIT %d", user.Id, user.DialogFromId, mode, limit)
	} else {
		rawSQL = fmt.Sprintf(baseSQL+" AND created_at >= NOW() - INTERVAL 1 HOUR ORDER BY created_at DESC LIMIT %d", user.Id, user.DialogFromId, mode, limit)
	}
	err := db.Raw("SELECT * FROM (" + rawSQL + ") AS subquery ORDER BY id ASC").Scan(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func GetUserLastMessage(user *models.User, db *gorm.DB) (*models.Message, error) {
	var message models.Message

	err := db.Where("user_id = ?", user.Id).
		Order("created_at DESC").
		Limit(1).
		Find(&message).Error

	if err != nil {
		return nil, err
	}

	return &message, nil
}
