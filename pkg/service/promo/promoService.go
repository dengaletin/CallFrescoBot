package promoService

import (
	"CallFrescoBot/pkg/models"
	promoRepository "CallFrescoBot/pkg/repositories/promo"
	"CallFrescoBot/pkg/utils"
	"errors"
)

func Create(campaignId uint64, user *models.User) error {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return errors.New("failed getting database connection")
	}

	return promoRepository.Create(campaignId, user, db)
}
