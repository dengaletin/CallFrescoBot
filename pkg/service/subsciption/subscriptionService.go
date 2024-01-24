package subsciptionService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	subscriptionRepository "CallFrescoBot/pkg/repositories/subscription"
	userRepository "CallFrescoBot/pkg/repositories/user"
	"CallFrescoBot/pkg/utils"
	"errors"
	"gorm.io/gorm"
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

func ResetSubscription(user *models.User) error {
	db, err := getDBConnection()
	if err != nil {
		return err
	}

	subscription, err := getSubscription(user, db)

	if subscription == nil && user.Mode == 2 || user.Dialog == 1 {
		resetErr := userRepository.ResetSubscription(user, db)
		if resetErr != nil {
			return resetErr
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
		subscription = &models.Subscription{Limit: consts.NoPlanLimit}

		return subscription, nil
	}

	return subscription, nil
}
