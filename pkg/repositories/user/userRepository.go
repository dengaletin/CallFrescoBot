package userRepository

import (
	"CallFrescoBot/pkg/models"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"time"
)

func FirstOrCreate(tgUser *tg.User, chatId int64, db *gorm.DB) (*models.User, error) {
	var user *models.User
	result := db.Where(models.User{TgId: tgUser.ID, Name: tgUser.UserName}).FirstOrCreate(&user)

	if result.RowsAffected == 0 && user.IsNew == true {
		db.Model(&user).Update("is_new", false)
	} else if result.RowsAffected == 0 {
		db.Model(&user).Update("last_login", time.Now())
	}

	db.Model(&user).Update("chat_id", chatId)

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

func SetMode(mode int64, user *models.User, db *gorm.DB) error {
	modes := []int64{0, 1}
	modeStatus := slices.Contains(modes, mode)

	if modeStatus == false {
		return errors.New("incorrect mode")
	}

	result := db.Model(&user).Update("mode", mode)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetMode(mode int64) (string, error) {
	switch mode {
	case 0:
		return "ChatGpt", nil
	case 1:
		return "Dalle3", nil
	default:
		return "", errors.New("unknown mode")
	}
}
