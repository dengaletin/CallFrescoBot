package payService

import (
	"CallFrescoBot/pkg/consts"
	"CallFrescoBot/pkg/models"
	invoiceRepository "CallFrescoBot/pkg/repositories/invoice"
	"CallFrescoBot/pkg/utils"
	"crypto/sha256"
	"encoding/hex"
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

func CreateInvoiceUrl(plan string, user *models.User) (string, error) {
	db, err := getDBConnection()
	if err != nil {
		return "", err
	}

	amount := resolveAmount(plan)
	currency := "USD"
	desc := resolveDescription(plan)
	lang := resolveLangById(user.Lang)
	coin := resolveCoin(plan)

	invoice, _ := invoiceRepository.InvoiceCreate(
		1,
		user.Id,
		amount,
		"USD",
		coin,
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

	baseURL := "https://aaio.io/merchant/pay"
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

func resolveAmount(plan string) float64 {
	switch plan {
	case "1":
		return 4.5
	case "2":
		return 9.9
	case "3":
		return 19.9
	default:
		return 0
	}
}

func resolveDescription(plan string) string {
	switch plan {
	case "1":
		return "Start"
	case "2":
		return "Pro"
	case "3":
		return "Boss"
	default:
		return ""
	}
}

func resolveCoin(plan string) int {
	switch plan {
	case "1":
		return 10
	case "2":
		return 20
	case "3":
		return 50
	default:
		return 0
	}
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
