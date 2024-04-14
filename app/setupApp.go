package app

import (
	paymentCallbackService "CallFrescoBot/pkg/service/paymentCallback"
	"CallFrescoBot/pkg/utils"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

func SetupApp() {
	log.Println("Initializing service")

	go setupServer()

	if err := createResources(); err != nil {
		log.Printf("Error occurred while setting up the app: %s", err)
	}
}

func setupServer() {
	http.HandleFunc("/payment-callback", paymentCallbackService.PaymentCallbackHandler)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createResources() error {
	if err := utils.CreateDBConnection(); err != nil {
		return errors.Wrap(err, "Failed to create database connection")
	}

	if err := utils.AutoMigrateDB(); err != nil {
		return errors.Wrap(err, "Failed to auto-migrate database")
	}

	if err := utils.SeedPlans(); err != nil {
		return errors.Wrap(err, "Failed to seed plans")
	}

	if err := utils.ClaudeUpdate(); err != nil {
		return errors.Wrap(err, "Failed to update subscriptions")
	}

	if err := utils.CreateBot(); err != nil {
		return errors.Wrap(err, "Failed to create bot")
	}

	return nil
}
