package userRefRepository

import (
	"CallFrescoBot/pkg/models"
	"gorm.io/gorm"
)

func Create(user *models.User, userRef *models.User, db *gorm.DB) (*models.UserRef, error) {
	ref := &models.UserRef{
		UserId1: user.Id,
		UserId2: userRef.Id,
	}

	result := db.Create(&ref)
	if result.Error != nil {
		return nil, result.Error
	}

	return ref, nil
}
