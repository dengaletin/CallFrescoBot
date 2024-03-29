package paymentCallbackService

import (
	"CallFrescoBot/pkg/consts"
	payService "CallFrescoBot/pkg/service/invoice"
	planService "CallFrescoBot/pkg/service/plan"
	subsciptionService "CallFrescoBot/pkg/service/subsciption"
	userService "CallFrescoBot/pkg/service/user"
	"CallFrescoBot/pkg/utils"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"strconv"
	"strings"
)

var secret2 = utils.GetEnvVar("AAIO_SECRET_2")

func calculateSignature(merchantID, amount, curr, secret, orderID string) string {
	data := strings.Join([]string{merchantID, amount, curr, secret, orderID}, ":")
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func PaymentCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "wrong request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}

	postedAmount := r.Form.Get("amount")
	postedCurrency := r.Form.Get("currency")
	orderId := r.Form.Get("order_id")
	u64, err := strconv.ParseUint(orderId, 10, 64)
	if err != nil {
		fmt.Println("Ошибка при конвертации:", err)
	}

	invoice, err := payService.GetInvoice(u64)
	if err != nil {
		fmt.Println("Ошибка при получении инвойса", err)
	}

	if invoice == nil {
		http.Error(w, "wrong invoice id", http.StatusBadRequest)
		return
	}

	stringAmount := strconv.FormatFloat(invoice.Amount, 'f', 2, 64)

	if invoice.Status != 0 {
		http.Error(w, "wrong invoice status", http.StatusBadRequest)
		return
	}

	if postedAmount != stringAmount {
		http.Error(w, "wrong amount", http.StatusBadRequest)
		return
	}

	if postedCurrency != invoice.Currency {
		http.Error(w, "wrong currency", http.StatusBadRequest)
		return
	}

	signature := calculateSignature(r.Form.Get("merchant_id"), postedAmount, postedCurrency, secret2, r.Form.Get("order_id"))
	if r.Form.Get("sign") != signature {
		http.Error(w, "wrong sign", http.StatusBadRequest)
		return
	}

	user, err := userService.GetUserById(invoice.UserId)
	invoice.Status = 1

	_, err = payService.UpdateInvoice(invoice)
	if err != nil {
		return
	}

	plan, err := planService.GetPlanById(*invoice.PlanId)
	if err != nil {
		return
	}

	_, err = subsciptionService.CreateWithPlan(user, plan)

	bot := utils.GetBot()
	msg := tg.NewMessage(user.TgId, utils.LocalizeSafe(consts.SubscriptionSuccess))
	_, err = bot.Send(msg)
	if err != nil {
		return
	}

	_, err = fmt.Fprintln(w, "OK")
	if err != nil {
		return
	}
}
