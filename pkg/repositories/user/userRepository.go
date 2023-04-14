package userRepository

import (
	"CallFrescoBot/pkg/models"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

func FirstOrCreate(tgUser *tg.User, db *gorm.DB) (*models.User, error) {
	var user *models.User
	err := db.Where(models.User{TgId: tgUser.ID, Name: tgUser.UserName}).Assign(models.User{LastLogin: time.Now()}).FirstOrCreate(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
