package subscriptionRepository

import (
	"CallFrescoBot/pkg/models"
	"gorm.io/gorm"
	"time"
)

func GetUserSubscription(user *models.User, db *gorm.DB) (*models.Subscription, error) {
	var subscription *models.Subscription

	result := db.Where("user_id = ? AND active_due > ?", user.Id, time.Now()).Find(&subscription)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected < 1 {
		subscription.Limit = 15
		return subscription, nil
	}

	return subscription, nil
}
