package subsciptionService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	subscriptionRepository "CallFrescoBot/pkg/repositories/subscription"
	"CallFrescoBot/pkg/utils"
	"errors"
)

func GetOrCreate(user *models.User, limit int, daysMultiplier int) (*models.Subscription, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	var subscription *models.Subscription

	subscription, err = subscriptionRepository.GetUserSubscription(user, db)
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
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, errors.New("error occurred while getting a DB connection from the connection pool")
	}

	subscription, err := subscriptionRepository.GetUserSubscription(user, db)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		subscription = &models.Subscription{Limit: consts.NoPlanLimit}
		return subscription, nil
	}

	return subscription, nil
}
