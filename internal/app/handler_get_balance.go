package app

import (
	"encoding/json"
	"net/http"
)

func (a *Application) HandlerGetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type result struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}

	user, err := a.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	account, err := a.Gophermart.GetAccount(r.Context(), *user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	current, _ := account.CurrentBalance.Float64()
	withdrawn, _ := account.WithdrawnBalance.Float64()
	responseBytes, err := json.Marshal(result{
		Current:   current,
		Withdrawn: withdrawn,
	})
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
