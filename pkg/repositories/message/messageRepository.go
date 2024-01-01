package messageRepository

import (
	"CallFrescoBot/pkg/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

func MessageCreate(userId uint64, message string, response string, db *gorm.DB) (*models.Message, error) {
	newMessage := models.Message{UserId: userId, Message: message, Response: response}
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

func GetMessagesByUser(user *models.User, limit int, db *gorm.DB) ([]models.Message, error) {
	var messages []models.Message

	err := db.Where("user_id = ?", user.Id).Order("created_at desc").Limit(limit).Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
