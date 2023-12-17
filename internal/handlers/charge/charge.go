package charge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"diplom.com/go-musthave-diploma-tpl/internal/authentication"
	"diplom.com/go-musthave-diploma-tpl/internal/dto"
	"diplom.com/go-musthave-diploma-tpl/internal/logger"
)

type ChargeHandler struct {
	Charge Charger
	Log    *logger.Logger
}

type Charger interface {
	ProcessOrder(ctx context.Context, order, userID string, sum float32) error
}

func NewChargeHandler(charge Charger, log *logger.Logger) *ChargeHandler {
	return &ChargeHandler{
		Charge: charge,
		Log:    log,
	}
}

func (c *ChargeHandler) ChargeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read data from request body", http.StatusBadRequest)
		return
	}

	var orderData dto.WithdrawOrder
	err = json.Unmarshal(body, &orderData)
	if err != nil {
		http.Error(w, "failed tro decode JSON data", http.StatusBadRequest)
		return
	}

	userID, _ := authentication.GetUserIDFromCtx(r.Context())

	err = c.Charge.ProcessOrder(context.Background(), orderData.Order, userID, orderData.Sum)
	if err != nil {
		c.Log.LogWarning("err when add withdraw to db: ", err)
		http.Error(w, "this order already exist, try another one", http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Order added.")

}
