package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/go-musthave-diploma/internal/gophermart"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func (a *Application) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := utils.Log

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := utils.ParseRequest(w, r, &req); err != nil {
		return
	}

	user, err := a.Gophermart.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		log.Infof("Error while registering new user: %v", err)
		var code int
		if errors.Is(err, storage.ErrDuplicateFound) {
			code = http.StatusConflict
		} else if errors.Is(err, gophermart.ErrEmptyLogin) || errors.Is(err, gophermart.ErrEmptyPassword) {
			code = http.StatusBadRequest
		} else {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
		return
	}

	cookie, err := a.createAuthCookie(*user)
	if err != nil {
		log.Infof("Error while issuing cookie: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
