package userRepository

import (
	"CallFrescoBot/pkg/models"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"time"
)

func FirstOrCreate(tgUser *tg.User, db *gorm.DB) (*models.User, error) {
	var user *models.User
	result := db.Where(models.User{TgId: tgUser.ID, Name: tgUser.UserName}).FirstOrCreate(&user)

	if result.RowsAffected == 0 && user.IsNew == true {
		db.Model(&user).Update("is_new", false)
	} else if result.RowsAffected == 0 {
		db.Model(&user).Update("last_login", time.Now())
	}

	return user, nil
}

func GerUserByTgId(tdId int64, db *gorm.DB) (*models.User, error) {
	var user *models.User
	err := db.Where(models.User{TgId: tdId}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
