package container

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

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
