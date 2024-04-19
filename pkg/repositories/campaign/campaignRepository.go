package campaignRepository

import (
	"CallFrescoBot/pkg/models"
	"errors"
	"gorm.io/gorm"
)

func Get(promoCode string, db *gorm.DB) (*models.Campaign, error) {
	var campaign *models.Campaign

	result := db.Where("code = ?", promoCode).First(&campaign)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return campaign, nil
}
