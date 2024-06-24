package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func (a *Application) HandlerGetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type row struct {
		Number     string              `json:"number"`
		Status     storage.OrderStatus `json:"status"`
		Accrual    *float64            `json:"accrual,omitempty"`
		UploadedAt time.Time           `json:"uploaded_at"`
	}

	user, err := a.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := a.Gophermart.GetOrders(r.Context(), *user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(*orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result := make([]row, len(*orders))
	for i, order := range *orders {
		var accrual *float64
		if order.Amount != nil {
			_accrual, _ := order.Amount.Float64()
			accrual = &_accrual
		}
		result[i] = row{
			Number:     order.OrderNumber,
			Status:     order.Status,
			Accrual:    accrual,
			UploadedAt: order.CreatedAt,
		}
	}

	responseBytes, err := json.Marshal(result)
	if err != nil {
		utils.Log.Errorf("Error during marshal: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
