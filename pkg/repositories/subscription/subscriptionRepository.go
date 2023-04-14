package subscriptionRepository

import (
	"CallFrescoBot/pkg/models"
	"gorm.io/gorm"
	"time"
)

func GetUserSubscriptionLimit(user *models.User, db *gorm.DB) int {
	var subscription *models.Subscription

	result := db.Where("user_id = ? AND active_due > ?", user.Id, time.Now()).Find(&subscription)

	if result.Error != nil || result.RowsAffected < 1 {
		return 15
	}

	return subscription.Limit
}
