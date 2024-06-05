package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	UserID           uuid.UUID       `db:"user_id"`
	CurrentBalance   decimal.Decimal `db:"current_balance"`
	WithdrawnBalance decimal.Decimal `db:"withdrawn_balance"`
}

func (s PgSQL) LoadAccount(ctx context.Context, user User) (*Account, error) {
	return loadAccount(ctx, s.Conn, user, false, false)
}

func saveAccount(ctx context.Context, querier pgxscan.Querier, account Account) error {
	_, err := querier.Query(
		ctx,
		`update public.account set current_balance = $1, withdrawn_balance = $2 where user_id = $3`,
		account.CurrentBalance, account.WithdrawnBalance, account.UserID,
	)
	if err != nil {
		return err
	}
	return nil
}

func loadAccount(ctx context.Context, querier pgxscan.Querier, user User, forUpdate bool, skipLocked bool) (*Account, error) {
	var result Account

	var suffixes []string
	if forUpdate {
		suffixes = append(suffixes, "for update")
	}
	if skipLocked {
		suffixes = append(suffixes, "skip locked")
	}
	query := fmt.Sprintf(`select * from public.account where user_id = $1 %s`, strings.Join(suffixes, " "))
	if err := pgxscan.Get(ctx, querier, &result, query, user.ID); err != nil {
		return nil, err
	}
	return &result, nil
}

func createAccount(ctx context.Context, querier pgxscan.Querier, user User) error {
	_, err := querier.Query(ctx, `insert into public.account (user_id) values ($1)`, user.ID)
	if err != nil {
		return err
	}
	return nil
}
