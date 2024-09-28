package utils

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	"CallFrescoBot/pkg/types"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var user string
var password string
var db string
var host string
var port string
var dbConn *gorm.DB

func init() {
	user = GetEnvVar("DB_USER")
	password = GetEnvVar("DB_PASSWORD")
	db = GetEnvVar("DB_NAME")
	host = GetEnvVar("DB_HOST")
	port = GetEnvVar("DB_PORT")
}

func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, db)
}

func CreateDBConnection() error {
	if dbConn != nil {
		CloseDBConnection(dbConn)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: GetDSN(),
	}), &gorm.Config{})

	if err != nil {
		log.Printf("Error occurred while connecting with the database")
	}

	sqlDB, err := db.DB()

	sqlDB.SetConnMaxLifetime(time.Hour)
	dbConn = db

	return err
}

func GetDatabaseConnection() (*gorm.DB, error) {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return dbConn, err
	}
	if err := sqlDB.Ping(); err != nil {
		return dbConn, err
	}
	return dbConn, nil
}

func CloseDBConnection(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		log.Printf("Error occurred while closing a DB connection")
	}
	defer sqlDB.Close()
}

func AutoMigrateDB() error {
	db, connErr := GetDatabaseConnection()
	if connErr != nil {
		return connErr
	}

	err := db.AutoMigrate(
		&models.User{},
		&models.Message{},
		&models.Subscription{},
		&models.UserRef{},
		&models.Invoice{},
		&models.Plan{},
		&models.Promo{},
		&models.Campaign{},
	)

	return err
}

func ClaudeUpdate() error {
	var subscriptions []models.Subscription
	err := dbConn.Find(&subscriptions).Error
	if err != nil {
		return err
	}

	for _, sub := range subscriptions {
		usage := make(map[string]interface{})
		if err := json.Unmarshal([]byte(sub.Usage), &usage); err != nil {
			return err
		}

		if _, exists := usage["claude"]; !exists {
			usage["claude"] = 0
		}
		if _, exists := usage["claude_context"]; !exists {
			usage["claude_context"] = 0
		}

		updatedUsageJSON, err := json.Marshal(usage)
		if err != nil {
			return err
		}

		err = dbConn.Model(&models.Subscription{}).Where("id = ?", sub.Id).Update("usage", updatedUsageJSON).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func SeedPlans() error {
	db, connErr := GetDatabaseConnection()
	if connErr != nil {
		return connErr
	}

	var count int64
	db.Model(&models.Plan{}).Count(&count)

	if count == 18 {
		limits := []types.Limit{
			{Gpt4OMiniLimit: 500, Gpt4OLimit: 300, Dalle3Limit: 30, Gpt4O1Limit: 0, ContextSupport: true},
		}

		configs := []types.Config{
			{limits[0], 0, 10},
		}

		var PlanNames = []string{
			consts.Plan19Name,
		}

		for index, config := range configs {
			configJSON, err := json.Marshal(config)
			if err != nil {
				return err
			}

			plan := models.Plan{
				Name:   PlanNames[index],
				Config: configJSON,
			}

			if err := db.Create(&plan).Error; err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}
