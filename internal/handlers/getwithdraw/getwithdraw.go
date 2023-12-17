package getwithdraw

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"diplom.com/go-musthave-diploma-tpl/internal/authentication"
	"diplom.com/go-musthave-diploma-tpl/internal/dto"
	"diplom.com/go-musthave-diploma-tpl/internal/logger"
)

type GetWithdrawHandler struct {
	GetWithDraw GetWithDrawer
	Log         *logger.Logger
}

var key = "UserID"

type GetWithDrawer interface {
	GetWithdrawals(ctx context.Context, userID string) ([]*dto.Billing, error)
}

func NewGetWithdrawHandler(getwd GetWithDrawer, log *logger.Logger) *GetWithdrawHandler {
	return &GetWithdrawHandler{
		GetWithDraw: getwd,
		Log:         log,
	}
}

func (g *GetWithdrawHandler) GetWithdrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only GET requests support!", http.StatusNotFound)
	}

	userID, ok := authentication.GetUserIDFromCtx(r.Context())
	if !ok {
		g.Log.LogWarning("could not get caller from context")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orders, err := g.GetWithDraw.GetWithdrawals(ctx, userID)
	if err != nil {
		g.Log.LogWarning("err when getting data from db:", err)
		fmt.Fprintf(w, "")
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
