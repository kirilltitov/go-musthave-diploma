package app

import (
	"errors"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/kirilltitov/go-musthave-diploma/internal/gophermart"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func (a *Application) HandlerWithdrawBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := utils.Log

	user, err := a.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var req struct {
		Order string          `json:"order"`
		Sum   decimal.Decimal `json:"sum"`
	}
	if err := utils.ParseRequest(w, r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.Gophermart.WithdrawBalanceFromAccount(r.Context(), *user, req.Sum, req.Order); err != nil {
		log.Infof("Error while withdrawing from account: %v", err)
		var code int
		if errors.Is(err, gophermart.ErrInvalidOrderNumber) {
			code = http.StatusUnprocessableEntity
		} else if errors.Is(err, storage.ErrInsufficientBalance) {
			code = http.StatusPaymentRequired
		} else {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
		return
	}

	w.WriteHeader(http.StatusOK)
}
