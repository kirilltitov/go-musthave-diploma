package storage

import (
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func (s PgSQL) InitDB(ctx context.Context) error {
	query := `SELECT count(*) cnt FROM pg_catalog.pg_tables where schemaname = 'public'`
	var result int
	if err := pgxscan.Get(ctx, s.Conn, &result, query); err != nil {
		return err
	}
	if result > 0 {
		utils.Log.Infof("Skipping migrations")
		return nil
	}

	migrations := []string{
		`create type public.order_status as enum ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')`,
		`create table public."user"
			(
				id         uuid        not null constraint user_pk primary key,
				login      varchar     not null constraint login_key unique,
				password   varchar     not null,
				created_at timestamptz not null
			)`,
		`create table public."order"
			(
				id           uuid           not null constraint order_pk primary key,
				order_number varchar        not null constraint order_number_key unique,
				user_id      uuid           not null,
				status       order_status   not null,
				amount       numeric(10, 2) not null,
				created_at   timestamptz    not null,
				updated_at   timestamptz
			)`,
		`create index order_created_at_index on "order" (created_at)`,
		`create index order_user_id_index on "order" (user_id)`,
		`create table public.account
			(
				user_id           uuid                     not null constraint account_pk primary key,
				current_balance   numeric(10, 2) default 0 not null,
				withdrawn_balance numeric(10, 2) default 0 not null
			)`,
		`create table public.withdrawal
			(
				id           uuid           not null constraint withdrawal_pk primary key,
				user_id      uuid           not null,
				order_number varchar        not null constraint withdrawal_order_number_ukey unique,
				amount       numeric(10, 2) not null,
				created_at   timestamptz    not null
			)`,
	}

	return WithVoidTransaction(ctx, s, func(tx pgx.Tx) error {
		for _, migration := range migrations {
			if _, err := tx.Exec(ctx, migration); err != nil {
				return err
			}
		}

		return tx.Commit(ctx)
	})
}
