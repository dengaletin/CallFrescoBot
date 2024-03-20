package planService

import (
	"CallFrescoBot/pkg/models"
	planRepository "CallFrescoBot/pkg/repositories/plan"
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

func GetPlanById(planId uint64) (*models.Plan, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	plan, err := planRepository.GetPlanById(planId, db)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func GetAllPlans() ([]*models.Plan, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	plan, err := planRepository.GetAllPlans(db)
	if err != nil {
		return nil, err
	}
	return plan, nil
}
