package gophermart

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

func (g Gophermart) WithdrawBalanceFromAccount(ctx context.Context, user storage.User, amount decimal.Decimal, order string) error {
	if err := validateOrderNumber(order); err != nil {
		return err
	}

	return g.container.Storage.WithdrawBalanceFromAccount(ctx, user, amount, order)
}
