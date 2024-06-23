package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

type Account struct {
	UserID           uuid.UUID       `db:"user_id"`
	CurrentBalance   decimal.Decimal `db:"current_balance"`
	WithdrawnBalance decimal.Decimal `db:"withdrawn_balance"`
}

func (s PgSQL) LoadAccount(ctx context.Context, user User) (*Account, error) {
	return WithTransaction(ctx, s, func(tx pgx.Tx) (*Account, error) {
		return loadAccount(ctx, tx, user, false, false)
	})
}

func (s PgSQL) ApplyProcessedOrder(ctx context.Context, user User, order Order, amount decimal.Decimal) error {
	return WithVoidTransaction(ctx, s, func(tx pgx.Tx) error {
		log := utils.Log.WithField("order_number", order.OrderNumber)

		account, err := loadAccount(ctx, tx, user, true, true)
		if err != nil {
			return err
		}

		account.CurrentBalance = account.CurrentBalance.Add(amount)
		if err := saveAccount(ctx, tx, *account); err != nil {
			return err
		}
		log.Debugf("Updated account")

		order.Amount = &amount
		if err := updateOrderAmount(ctx, tx, order); err != nil {
			return err
		}
		log.Debugf("Updated order amount")
		if err := updateOrderStatus(ctx, tx, order, StatusProcessed, []OrderStatus{StatusNew, StatusProcessing, StatusProcessed}); err != nil {
			return err
		}
		log.Debugf("Updated order status")

		return tx.Commit(ctx)
	})
}

func saveAccount(ctx context.Context, tx pgx.Tx, account Account) error {
	utils.Log.Debugf("About to update account %s with current balance = %f and withdrawn balance = %f", account.UserID, account.CurrentBalance.InexactFloat64(), account.WithdrawnBalance.InexactFloat64())

	res, err := tx.Exec(
		ctx,
		`update public.account set current_balance = $1, withdrawn_balance = $2 where user_id = $3`,
		account.CurrentBalance.InexactFloat64(), account.WithdrawnBalance.InexactFloat64(), account.UserID,
	)
	if err != nil {
		return err
	}

	utils.Log.Debugf("Account update result: %+v", res)

	return nil
}

func loadAccount(ctx context.Context, tx pgx.Tx, user User, forUpdate bool, skipLocked bool) (*Account, error) {
	var result Account

	var suffixes []string
	if forUpdate {
		suffixes = append(suffixes, "for update")
	}
	if skipLocked {
		suffixes = append(suffixes, "skip locked")
	}
	query := fmt.Sprintf(`select * from public.account where user_id = $1 %s`, strings.Join(suffixes, " "))
	if err := pgxscan.Get(ctx, tx, &result, query, user.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &result, nil
}

func createAccount(ctx context.Context, tx pgx.Tx, user User) error {
	_, err := tx.Exec(ctx, `insert into public.account (user_id) values ($1)`, user.ID)
	if err != nil {
		return err
	}
	return nil
}
