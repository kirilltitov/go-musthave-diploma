package app

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

func (a *Application) HandlerGetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type row struct {
		Number     string              `json:"number"`
		Status     storage.OrderStatus `json:"status"`
		Accrual    float64             `json:"accrual"`
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

	var result []row
	for _, order := range *orders {
		accrual, _ := order.Amount.Float64()
		result = append(result, row{
			Number:     order.OrderNumber,
			Status:     order.Status,
			Accrual:    accrual,
			UploadedAt: order.CreatedAt,
		})
	}

	responseBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
