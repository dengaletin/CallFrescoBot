package payService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	invoiceRepository "CallFrescoBot/pkg/repositories/invoice"
	"CallFrescoBot/pkg/types"
	"CallFrescoBot/pkg/utils"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"strconv"
)

var merchantId = utils.GetEnvVar("AAIO_MERCHANT_ID")
var secret = utils.GetEnvVar("AAIO_SECRET_1")

func getDBConnection() (*gorm.DB, error) {
	db, err := utils.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func UpdateInvoice(invoice *models.Invoice) (*models.Invoice, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	invoice, err = invoiceRepository.UpdateInvoice(invoice, db)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func GetInvoice(invoiceId uint64) (*models.Invoice, error) {
	db, err := getDBConnection()
	if err != nil {
		return nil, err
	}

	invoice, err := invoiceRepository.GetInvoice(invoiceId, db)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func CreateInvoiceUrl(plan *models.Plan, user *models.User) (string, error) {
	db, err := getDBConnection()
	if err != nil {
		return "", err
	}

	var price float64
	var currency string
	var language string
	var config types.Config
	if err := json.Unmarshal(plan.Config, &config); err != nil {
		return "", err
	}

	switch user.Lang {
	case consts.LangEn:
		price = config.PriceEn
		language = "en"
		currency = "USD"
	case consts.LangRu:
		price = config.PriceRu
		language = "ru"
		currency = "RUB"
	default:
		return "", errors.New("unknown language")
	}

	amount := price

	desc := plan.Name
	lang := language
	planId := plan.Id

	invoice, _ := invoiceRepository.InvoiceCreate(
		1,
		user.Id,
		amount,
		currency,
		0,
		&planId,
		0,
		db,
	)

	orderID := invoice.Id

	data := []string{
		merchantId,
		strconv.FormatFloat(amount, 'f', 2, 64),
		currency,
		secret,
		strconv.FormatUint(orderID, 10),
	}
	hash := sha256.New()
	hash.Write([]byte(join(data, ":")))
	signature := hex.EncodeToString(hash.Sum(nil))

	baseURL := "https://aaio.so/merchant/pay"
	params := url.Values{}
	params.Set("merchant_id", merchantId)
	params.Set("amount", fmt.Sprintf("%.2f", amount))
	params.Set("currency", currency)
	params.Set("order_id", strconv.FormatUint(orderID, 10))
	params.Set("sign", signature)
	params.Set("desc", desc)
	params.Set("lang", lang)
	fmt.Println(baseURL + "?" + params.Encode())

	return baseURL + "?" + params.Encode(), nil
}

func join(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	if len(slice) == 1 {
		return slice[0]
	}
	str := slice[0]
	for i := 1; i < len(slice); i++ {
		str += sep + slice[i]
	}
	return str
}

func resolveLangById(id int64) string {
	switch id {
	case consts.LangEn:
		return "en"
	case consts.LangRu:
		return "ru"
	default:
		return "en"
	}
}
