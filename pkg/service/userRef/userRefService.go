package userRefService

import (
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/repositories/userRef"
	"CallFrescoBot/pkg/utils"
	"errors"
)

func Create(user, refUser *models.User) (*models.UserRef, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("failed getting database connection")
	}

	return userRefRepository.Create(user, refUser, db)
}
