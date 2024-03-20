package subsciptionService

import (
	"CallFrescoBot/pkg/models"
	subscriptionRepository "CallFrescoBot/pkg/repositories/subscription"
	userRepository "CallFrescoBot/pkg/repositories/user"
	"CallFrescoBot/pkg/types"
	"CallFrescoBot/pkg/utils"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

func getDBConnection() (*gorm.DB, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}
	return db, nil
}

func GetUserSubscription(user *models.User) (*models.Subscription, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	subscription, err := getSubscription(user, db)

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func UpdateUserSubscription(subscription *models.Subscription) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}

	_, err = subscriptionRepository.UpdateSubscription(subscription, db)
	if err != nil {
		return err
	}

	return nil
}

func ResetSubscription(user *models.User) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}

	subscription, err := getSubscription(user, db)

	if (subscription == nil || subscription.PlanId == nil) && (user.Mode != 0 || user.Dialog != 0) {
		resetErr := userRepository.ResetSubscription(user, db)
		if resetErr != nil {
			return resetErr
		}
	}

	if subscription == nil || subscription.PlanId == nil {
		return nil
	}

	currentDate := time.Now().Truncate(24 * time.Hour)
	dueDate := subscription.ActiveDue.Truncate(24 * time.Hour)
	refreshDate := subscription.RefreshDate.Truncate(24 * time.Hour)

	if currentDate.After(dueDate) {
		return nil
	}

	if currentDate.Equal(refreshDate) {
		newUsage := types.Usage{}
		newUsageJson, err := json.Marshal(newUsage)
		if err != nil {
			return err
		}

		subscription.Usage = newUsageJson
		subscription.RefreshDate = currentDate.AddDate(0, 0, 30)

		if _, err := subscriptionRepository.UpdateSubscription(subscription, db); err != nil {
			return err
		}
	}

	return nil
}

func getSubscription(user *models.User, db *gorm.DB) (*models.Subscription, error) {
	subscription, err := subscriptionRepository.GetUserSubscription(user, db)
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func Create(user *models.User, limit int) (*models.Subscription, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	var subscription *models.Subscription
	subscription, err = subscriptionRepository.CreateSubscription(user, limit, 30, db)

	return subscription, nil
}

func CreateWithPlan(user *models.User, plan *models.Plan) (*models.Subscription, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	var subscription *models.Subscription
	subscription, err = subscriptionRepository.CreateSubscriptionWithPlan(user, plan, 30, db)

	return subscription, nil
}

func GetOrCreate(user *models.User, limit int, daysMultiplier int) (*models.Subscription, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	var subscription *models.Subscription

	subscription, err = getSubscription(user, db)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		subscription, err = subscriptionRepository.CreateSubscription(user, limit, daysMultiplier, db)
	} else {
		if subscription.Limit > limit {
			return nil, errors.New("subscription is too cool")
		}

		subscription.ActiveDue = subscription.ActiveDue.AddDate(0, 0, daysMultiplier)
		subscription, err = subscriptionRepository.UpdateSubscription(subscription, db)
	}

	return subscription, nil
}

func GetUserSubscriptionWithNoPlanLimit(user *models.User) (*models.Subscription, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	subscription, _ := getSubscription(user, db)
	if subscription == nil {
		return nil, nil
	}

	return subscription, nil
}
