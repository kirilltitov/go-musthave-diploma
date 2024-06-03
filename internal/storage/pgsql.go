package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

type PgSQL struct {
	Conn *pgxpool.Pool
}

func New(ctx context.Context, DSN string) (*PgSQL, error) {
	conf, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}
	conf.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	utils.Log.Infof("Connected to PgSQL with DSN %s", DSN)

	return &PgSQL{Conn: pool}, nil
}

func WithTransaction[T any](ctx context.Context, pg PgSQL, f func(pgx.Tx) (*T, error)) (*T, error) {
	transaction, err := pg.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	result, err := f(transaction)
	if err != nil {
		transaction.Rollback(ctx)
		return nil, err
	}

	return result, nil
}
