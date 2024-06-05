package app

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/container"
	"github.com/kirilltitov/go-musthave-diploma/internal/gophermart"
)

type Application struct {
	Config     config.Config
	Container  *container.Container
	Gophermart gophermart.Gophermart
}

func New(ctx context.Context, cfg config.Config, cnt *container.Container) (*Application, error) {
	if cnt == nil {
		_cnt, err := container.New(ctx, cfg)
		if err != nil {
			return nil, err
		}
		cnt = _cnt
	}

	return &Application{
		Config:     cfg,
		Container:  cnt,
		Gophermart: gophermart.New(cfg, cnt),
	}, nil
}

func (a *Application) Run() error {
	r := a.createRouter()

	r.Use(
		middleware.RequestID,
		middleware.Compress(5),
		middleware.Logger,
		middleware.Recoverer,
	)

	return http.ListenAndServe(a.Config.ServerAddress, r)
}

func (a *Application) createRouter() chi.Router {
	router := chi.NewRouter()

	router.Post("/api/user/register", a.HandlerRegister)
	router.Post("/api/user/login", a.HandlerLogin)
	//router.Post("/api/user/orders", a.HandlerCreateOrder)
	router.Get("/api/user/orders", a.HandlerGetOrders)
	router.Get("/api/user/balance", a.HandlerGetBalance)
	router.Post("/api/user/balance/withdraw", a.HandlerWithdrawBalance)
	router.Post("/api/user/withdrawals", a.HandlerGetWithdrawals)

	return router
}
