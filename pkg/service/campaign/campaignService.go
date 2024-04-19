package campaignService

import (
	"CallFrescoBot/pkg/models"
	campaignRepository "CallFrescoBot/pkg/repositories/campaign"
	"CallFrescoBot/pkg/utils"
	"errors"
)

func Get(promoCode string) (*models.Campaign, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("failed getting database connection")
	}

	return campaignRepository.Get(promoCode, db)
}
