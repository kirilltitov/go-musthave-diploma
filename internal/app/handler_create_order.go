package app

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/kirilltitov/go-musthave-diploma/internal/gophermart"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func (a *Application) HandlerCreateOrder(w http.ResponseWriter, r *http.Request) {
	user, err := a.authenticate(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var buf bytes.Buffer
	defer r.Body.Close()

	if _, err := buf.ReadFrom(r.Body); err != nil {
		utils.Log.Errorf("Could not get body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	orderNumber := buf.String()
	if orderNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := a.Gophermart.CreateOrder(r.Context(), *user, orderNumber); err != nil {
		utils.Log.Errorf("Error while creating order: %v", err)
		switch {
		case errors.Is(err, gophermart.ErrInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		case errors.Is(err, gophermart.ErrOrderAlreadyUploaded):
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, gophermart.ErrNotYourOrder):
			w.WriteHeader(http.StatusConflict)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
