package gophermart

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

// GetWithdrawals Returns withdrawals for given user
func (g Gophermart) GetWithdrawals(ctx context.Context, user storage.User) (*[]storage.Withdrawal, error) {
	return g.container.Storage.LoadWithdrawals(ctx, user)
}
