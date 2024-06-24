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

	return http.ListenAndServe(a.Config.ServerAddress, r)
}

func (a *Application) createRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.Compress(5),
		middleware.Logger,
		middleware.Recoverer,
	)

	r.Post("/api/user/register", a.HandlerRegister)
	r.Post("/api/user/login", a.HandlerLogin)
	r.Post("/api/user/orders", a.HandlerCreateOrder)
	r.Get("/api/user/orders", a.HandlerGetOrders)
	r.Get("/api/user/balance", a.HandlerGetBalance)
	r.Post("/api/user/balance/withdraw", a.HandlerWithdrawBalance)
	r.Get("/api/user/withdrawals", a.HandlerGetWithdrawals)

	return r
}
