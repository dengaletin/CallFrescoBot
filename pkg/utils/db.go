package utils

import (
	"CallFrescoBot/pkg/models"
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

	err := db.AutoMigrate(&models.User{}, &models.Message{}, &models.Subscription{}, &models.UserRef{})
	return err
}
