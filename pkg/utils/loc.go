package utils

import (
	"CallFrescoBot/pkg/consts"
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"path/filepath"
)

var localizer *i18n.Localizer

func InitBundle(langID int64) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	messageFiles := []string{"en.json", "ru.json"}
	for _, f := range messageFiles {
		filePath := filepath.Join("pkg/messages", f)
		_, err := bundle.LoadMessageFile(filePath)
		if err != nil {
			log.Fatalf("failed to load message file '%s': %v", filePath, err)
		}
	}

	lang := resolveLangById(langID)
	localizer = i18n.NewLocalizer(bundle, lang)
}

func LocalizeSafe(messageID string) string {
	localized, err := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: messageID,
		},
	})
	if err != nil {
		log.Printf("Error localizing message '%s': %v", messageID, err)
		return messageID
	}

	return localized
}

func resolveLangById(id int64) string {
	switch id {
	case consts.LangEn:
		return consts.LangEnName
	case consts.LangRu:
		return consts.LangRuName
	default:
		return consts.LangEnName
	}
}
