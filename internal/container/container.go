package container

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/accrual"
	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

type Container struct {
	Storage storage.Storage
	Accrual accrual.Accrual
}

func New(ctx context.Context, cfg config.Config) (*Container, error) {
	s, err := newPgSQLStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Container{
		Storage: *s,
		Accrual: accrual.NewExternalAccrual(
			accrual.ExternalAccrualConfig{
				Address: cfg.AccrualSystemAddress,
				Timeout: 1,
				Retries: 5,
			},
		),
	}, nil
}

func newPgSQLStorage(ctx context.Context, cfg config.Config) (*storage.PgSQL, error) {
	s, err := storage.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	if err := s.InitDB(ctx); err != nil {
		return nil, err
	}

	return s, nil
}
