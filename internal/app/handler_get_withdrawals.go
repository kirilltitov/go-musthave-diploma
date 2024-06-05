package app

import (
	"encoding/json"
	"net/http"
	"time"
)

func (a *Application) HandlerGetWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type row struct {
		Order       string    `json:"order"`
		Sum         float64   `json:"sum"`
		ProcessedAt time.Time `json:"processed_at"`
	}

	user, err := a.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdrawals, err := a.Gophermart.GetWithdrawals(r.Context(), *user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(*withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var result []row
	for _, withdrawal := range *withdrawals {
		sum, _ := withdrawal.Amount.Float64()
		result = append(result, row{
			Order:       withdrawal.OrderNumber,
			Sum:         sum,
			ProcessedAt: withdrawal.CreatedAt,
		})
	}

	responseBytes, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
