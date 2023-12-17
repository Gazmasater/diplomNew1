package addorder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"diplom.com/go-musthave-diploma-tpl/internal/authentication"
	"diplom.com/go-musthave-diploma-tpl/internal/logger"
	"github.com/ShiraazMoollatjie/goluhn"
)

var (
	ErrOrderExists   = errors.New("this order already exists, please try another one")
	ErrDifferentUser = errors.New("another user has already placed an order with this number")
)

type AddOrderHandler struct {
	addOrder    AddOrder
	RedisClient *redis.RedisClient
	Log         *logger.Logger
}

type AddOrder interface {
	OrderSetter(ctx context.Context, userID, orderNumber string) error
}

func NewPutOrderNumberHandler(order AddOrder, redis *redis.RedisClient, log *logger.Logger) *AddOrderHandler {
	return &AddOrderHandler{
		addOrder:    order,
		RedisClient: redis,
		Log:         log,
	}
}

func (a *AddOrderHandler) AddOrderNumberHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST requests support!", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read data from request body", http.StatusBadRequest)
		return
	}

	orderNumber := string(body)
	err = goluhn.Validate(orderNumber)
	if err != nil {
		http.Error(w, "invalid format of order number", http.StatusUnprocessableEntity)
		return
	}

	userID, ok := authentication.GetUserIDFromCtx(r.Context())
	if !ok {
		a.Log.LogWarning("could not get caller from context")
	}

	err = a.addOrder.OrderSetter(context.Background(), userID, orderNumber)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderExists):
			http.Error(w, "this order already exists", http.StatusOK)
		case errors.Is(err, ErrDifferentUser):
			http.Error(w, "another user has already placed an order with this number", http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Order added.")
}
