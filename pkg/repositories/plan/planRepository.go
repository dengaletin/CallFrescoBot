package planRepository

import (
	"CallFrescoBot/pkg/models"
	"errors"
	"gorm.io/gorm"
)

func GetPlanById(planId uint64, db *gorm.DB) (*models.Plan, error) {
	var plan *models.Plan

	result := db.Where("id = ?", planId).First(&plan)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return plan, nil
}

func GetAllPlans(db *gorm.DB) ([]*models.Plan, error) {
	var plans []*models.Plan

	result := db.Find(&plans)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}

	return plans, nil
}
