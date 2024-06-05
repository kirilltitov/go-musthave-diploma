package storage

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

type Withdrawal struct {
	ID          uuid.UUID       `db:"id"`
	UserID      uuid.UUID       `db:"user_id"`
	OrderNumber string          `db:"order_number"`
	Amount      decimal.Decimal `db:"amount"`
	CreatedAt   time.Time       `db:"created_at"`
}

func (s PgSQL) LoadWithdrawals(ctx context.Context, user User) (*[]Withdrawal, error) {
	var rows []Withdrawal

	err := pgxscan.Select(
		ctx,
		s.Conn,
		&rows,
		`select * from public.withdrawal where user_id = $1 order by created_at asc`,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

func (s PgSQL) WithdrawBalanceFromAccount(ctx context.Context, user User, amount decimal.Decimal, order string) error {
	return WithVoidTransaction(ctx, s, func(tx pgx.Tx) error {
		account, err := loadAccount(ctx, tx, user, true, true)
		if err != nil {
			return err
		}

		if account.CurrentBalance.LessThan(amount) {
			return ErrInsufficientBalance
		}

		account.CurrentBalance.Sub(amount)
		account.WithdrawnBalance.Add(amount)

		if err := saveAccount(ctx, tx, *account); err != nil {
			return err
		}

		withdrawal := Withdrawal{
			ID:          utils.NewUUID6(),
			UserID:      user.ID,
			OrderNumber: order,
			Amount:      amount,
			CreatedAt:   time.Now(),
		}
		if err := createWithdrawal(ctx, tx, withdrawal); err != nil {
			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return err
		}

		return nil
	})
}

func createWithdrawal(ctx context.Context, querier pgxscan.Querier, withdrawal Withdrawal) error {
	_, err := querier.Query(
		ctx,
		`insert into public.withdrawal (id, user_id, order_number, amount, created_at) values ($1, $2, $3, $4, $5)`,
		withdrawal.ID,
		withdrawal.UserID,
		withdrawal.OrderNumber,
		withdrawal.Amount,
		withdrawal.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
