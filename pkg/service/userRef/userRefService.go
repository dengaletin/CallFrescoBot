package userRefService

import (
	"CallFrescoBot/pkg/models"
	userRefRepository "CallFrescoBot/pkg/repositories/userRef"
	"CallFrescoBot/pkg/utils"
	"errors"
)

func Create(user *models.User, refUser *models.User) (*models.UserRef, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	userRef, err := userRefRepository.Create(user, refUser, db)
	if err != nil {
		return nil, err
	}

	return userRef, nil
}
