package promoRepository

import (
	"CallFrescoBot/pkg/models"
	"gorm.io/gorm"
)

func Create(campaignId uint64, user *models.User, db *gorm.DB) error {
	ref := &models.Promo{
		UserId:     user.Id,
		CampaignId: campaignId,
	}

	result := db.Create(&ref)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
