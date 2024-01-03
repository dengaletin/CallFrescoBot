package app

import (
	"CallFrescoBot/pkg/utils"
	"github.com/pkg/errors"
	"log"
)

func SetupApp() {
	log.Println("Initializing service")

	if err := createResources(); err != nil {
		log.Printf("Error occurred while setting up the app: %s", err)
	}
}

func createResources() error {
	if err := utils.CreateDBConnection(); err != nil {
		return errors.Wrap(err, "Failed to create database connection")
	}

	if err := utils.AutoMigrateDB(); err != nil {
		return errors.Wrap(err, "Failed to auto-migrate database")
	}

	if err := utils.CreateBot(); err != nil {
		return errors.Wrap(err, "Failed to create bot")
	}

	return nil
}
