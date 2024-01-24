package invoiceRepository

import (
	"CallFrescoBot/pkg/models"
	"errors"
	"gorm.io/gorm"
)

func GetInvoice(invoiceId uint64, db *gorm.DB) (*models.Invoice, error) {
	var invoice *models.Invoice
	err := db.Where(models.Invoice{Id: invoiceId}).First(&invoice).Error
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func InvoiceCreate(
	paymentMethodId uint64,
	userId uint64,
	amount float64,
	currency string,
	coin int,
	status int64,
	db *gorm.DB,
) (*models.Invoice, error) {
	newInvoice := models.Invoice{
		PaymentMethodId: paymentMethodId,
		UserId:          userId,
		Amount:          amount,
		Currency:        currency,
		Coin:            coin,
		Status:          status,
	}
	result := db.Create(&newInvoice)

	if result.Error != nil && result.RowsAffected != 1 {
		return nil, errors.New("error occurred while creating a new invoice")
	}

	return &newInvoice, nil
}

func UpdateInvoice(invoice *models.Invoice, db *gorm.DB) (*models.Invoice, error) {
	result := db.Save(invoice)
	if result.Error != nil {
		return nil, result.Error
	}

	return invoice, nil
}
