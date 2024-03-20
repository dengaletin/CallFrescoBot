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
	)

	return err
}

func SeedPlans() error {
	db, connErr := GetDatabaseConnection()
	if connErr != nil {
		return connErr
	}

	var count int64
	db.Model(&models.Plan{}).Count(&count)

	if count > 0 {
		return nil
	}

	limits := []types.Limit{
		{Gpt35Limit: 100, Gpt4Limit: 0, Dalle3Limit: 0, ContextSupport: false},
		{Gpt35Limit: 500, Gpt4Limit: 0, Dalle3Limit: 0, ContextSupport: false},
		{Gpt35Limit: 100, Gpt4Limit: 0, Dalle3Limit: 0, ContextSupport: true},
		{Gpt35Limit: 1000, Gpt4Limit: 0, Dalle3Limit: 0, ContextSupport: false},
		{Gpt35Limit: 500, Gpt4Limit: 0, Dalle3Limit: 0, ContextSupport: true},
		{Gpt35Limit: 100, Gpt4Limit: 100, Dalle3Limit: 0, ContextSupport: false},
		{Gpt35Limit: 100, Gpt4Limit: 100, Dalle3Limit: 10, ContextSupport: false},
		{Gpt35Limit: 1000, Gpt4Limit: 0, Dalle3Limit: 0, ContextSupport: true},
		{Gpt35Limit: 100, Gpt4Limit: 100, Dalle3Limit: 0, ContextSupport: true},
		{Gpt35Limit: 100, Gpt4Limit: 100, Dalle3Limit: 10, ContextSupport: true},
		{Gpt35Limit: 500, Gpt4Limit: 500, Dalle3Limit: 0, ContextSupport: false},
		{Gpt35Limit: 500, Gpt4Limit: 500, Dalle3Limit: 50, ContextSupport: false},
		{Gpt35Limit: 1000, Gpt4Limit: 1000, Dalle3Limit: 0, ContextSupport: false},
		{Gpt35Limit: 1000, Gpt4Limit: 1000, Dalle3Limit: 100, ContextSupport: false},
	}

	configs := []types.Config{
		{limits[0], 95, 1},
		{limits[1], 185, 2},
		{limits[2], 275, 3},
		{limits[3], 320, 3.50},
		{limits[4], 549, 6},
		{limits[5], 730, 8},
		{limits[6], 825, 9},
		{limits[7], 959, 10.50},
		{limits[8], 2195, 24},
		{limits[9], 2466, 27},
		{limits[10], 3199, 35},
		{limits[11], 3654, 40},
		{limits[12], 5938, 65},
		{limits[13], 6852, 75},
	}

	var PlanNames = []string{
		consts.Plan1Name,
		consts.Plan2Name,
		consts.Plan3Name,
		consts.Plan4Name,
		consts.Plan5Name,
		consts.Plan6Name,
		consts.Plan7Name,
		consts.Plan8Name,
		consts.Plan9Name,
		consts.Plan10Name,
		consts.Plan11Name,
		consts.Plan12Name,
		consts.Plan13Name,
		consts.Plan14Name,
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
