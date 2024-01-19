package utils

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"encoding/json"
	"errors"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"path/filepath"
)

var loc *i18n.Localizer

func CreateLoc(user *models.User) error {
	lang, err := resolveLangById(user.Lang)
	if err != nil {
		return err
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err = bundle.LoadMessageFile(filepath.Join("pkg/messages", "en.json"))
	if err != nil {
		return err
	}
	_, err = bundle.LoadMessageFile(filepath.Join("pkg/messages", "ru.json"))
	if err != nil {
		return err
	}

	loc = i18n.NewLocalizer(bundle, lang)

	return err
}

func LocalizeSafe(messageID string) string {
	localizer := getLoc()

	localized, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		log.Printf("Error localizing message '%s': %v", messageID, err)

		return messageID
	}

	return localized
}

func getLoc() *i18n.Localizer {
	return loc
}

func resolveLangById(id int64) (string, error) {
	switch id {
	case consts.LangEn:
		return "en", nil
	case consts.LangRu:
		return "ru", nil
	default:
		return "", errors.New("unsupported language")
	}
}
