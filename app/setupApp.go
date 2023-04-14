package app

import (
	"CallFrescoBot/pkg/utils"
	"log"
)

func SetupApp() {
	log.Printf("Initializing service")

	if err := utils.CreateDBConnection(); err != nil {
		log.Printf("Error occurred while creating the database connection")
		return
	}

	err := utils.AutoMigrateDB()
	if err != nil {
		log.Printf("Error occurred while auto migrating database")
		return
	}

	err = utils.CreateBot()
	if err != nil {
		log.Printf("Error occurred while bot creating")
		return
	}
}
