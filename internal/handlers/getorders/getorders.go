package getorders

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"diplom.com/go-musthave-diploma-tpl/internal/authentication"
	"diplom.com/go-musthave-diploma-tpl/internal/dto"
	"diplom.com/go-musthave-diploma-tpl/internal/interrogator"
	"diplom.com/go-musthave-diploma-tpl/internal/logger"
)

type GetOrdersHandler struct {
	log          *logger.Logger
	OrdersGetter OrdersGetter
	interrogator *interrogator.Interrogator
}

type OrdersGetter interface {
	GetOrders(ctx context.Context, userID string) ([]dto.FullOrder, error)
}

func NewGetOrdersHandler(ord OrdersGetter, log *logger.Logger, interrogator *interrogator.Interrogator) *GetOrdersHandler {
	return &GetOrdersHandler{
		log:          log,
		OrdersGetter: ord,
		interrogator: interrogator,
	}
}

func (g *GetOrdersHandler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests support!", http.StatusNotFound)
	}

	userID, ok := authentication.GetUserIDFromCtx(r.Context())
	if !ok {
		g.log.LogWarning("could not get caller from context")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)

	defer cancel()

	orders, err := g.OrdersGetter.GetOrders(ctx, userID)
	if err != nil {
		g.log.LogWarning("err when getting orders list:", err)
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	}

	jsonData, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
