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
	var user models.User
	result := db.Where("tg_id = ?", tgUser.ID).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newUser := models.User{
			TgId:  tgUser.ID,
			Name:  tgUser.UserName,
			IsNew: true,
		}
		if err := db.Create(&newUser).Error; err != nil {
			return nil, err
		}

		return &newUser, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	updates := map[string]interface{}{}
	if user.IsNew {
		updates["is_new"] = false
	}
	if user.Name != tgUser.UserName {
		updates["name"] = tgUser.UserName
	}
	updates["last_login"] = time.Now()
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GerUserByTgId(tdId int64, db *gorm.DB) (*models.User, error) {
	var user *models.User
	err := db.Where(models.User{TgId: tdId}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserById(id uint64, db *gorm.DB) (*models.User, error) {
	var user *models.User
	err := db.Where(models.User{Id: id}).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func SetMode(mode int64, user *models.User, db *gorm.DB) error {
	modes := []int64{consts.Gpt4oMiniMode, consts.DalleMode, consts.Gpt4oMode, consts.Gpt4o1Mode}
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
	languages := []int64{consts.LangEn, consts.LangRu}
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
	allowedValue := []int64{consts.DialogModeOff, consts.DialogModeOn}
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
	case consts.Gpt4oMiniMode:
		return utils.LocalizeSafe(consts.ModeGpt4oMini), nil
	case consts.DalleMode:
		return utils.LocalizeSafe(consts.ModeDalle3), nil
	case consts.Gpt4oMode:
		return utils.LocalizeSafe(consts.ModeGpt4), nil
	case consts.Gpt4o1Mode:
		return utils.LocalizeSafe(consts.ModeGpt4o1), nil
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
		"mode":   consts.Gpt4oMiniMode,
		"dialog": consts.DialogModeOff,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
