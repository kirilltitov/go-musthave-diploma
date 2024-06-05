package gophermart

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

func (g Gophermart) GetWithdrawals(ctx context.Context, user storage.User) (*[]storage.Withdrawal, error) {
	return g.container.Storage.LoadWithdrawals(ctx, user)
}
