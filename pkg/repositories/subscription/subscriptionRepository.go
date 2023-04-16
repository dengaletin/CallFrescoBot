package subscriptionRepository

import (
	"CallFrescoBot/pkg/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

func GetUserSubscription(user *models.User, db *gorm.DB) (*models.Subscription, error) {
	var subscription *models.Subscription

	result := db.Where("user_id = ? AND active_due > ?", user.Id, time.Now()).Last(&subscription)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return subscription, nil
}

func CreateSubscription(user *models.User, limit int, multiplierDays int, db *gorm.DB) (*models.Subscription, error) {
	subscription := &models.Subscription{
		UserId:    user.Id,
		Limit:     limit,
		ActiveDue: time.Now().AddDate(0, 0, multiplierDays),
	}

	result := db.Create(&subscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}

func UpdateSubscription(subscription *models.Subscription, db *gorm.DB) (*models.Subscription, error) {
	result := db.Save(subscription)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscription, nil
}
