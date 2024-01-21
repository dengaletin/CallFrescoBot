package userRepository

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/utils"
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slices"
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

func SetMode(mode int64, user *models.User, db *gorm.DB) error {
	modes := []int64{0, 1, 2}
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

func SetLanguage(language int64, user *models.User, db *gorm.DB) error {
	languages := []int64{1, 2}
	languageStatus := slices.Contains(languages, language)

	if languageStatus == false {
		return errors.New("incorrect language")
	}

	result := db.Model(&user).Update("lang", language)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SetDialogStatus(dialogStatus int64, user *models.User, db *gorm.DB) error {
	allowedValue := []int64{0, 1}
	dialogStatusValue := slices.Contains(allowedValue, dialogStatus)

	if dialogStatusValue == false {
		return errors.New("incorrect value")
	}

	result := db.Model(&user).Update("dialog", dialogStatus)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetMode(mode int64) (string, error) {
	switch mode {
	case 0:
		return utils.LocalizeSafe(consts.ModeGpt35), nil
	case 1:
		return utils.LocalizeSafe(consts.ModeDalle3), nil
	case 2:
		return utils.LocalizeSafe(consts.ModeGpt4), nil
	default:
		return "", errors.New("unknown mode")
	}
}

func SetDialogFromId(messageId uint64, user *models.User, db *gorm.DB) error {
	result := db.Model(&user).Update("dialog_from_id", messageId)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ResetSubscription(user *models.User, db *gorm.DB) error {
	result := db.Model(user).Updates(map[string]interface{}{
		"mode":   0,
		"dialog": 0,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
