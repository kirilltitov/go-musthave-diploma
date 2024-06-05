package gophermart

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

func (g Gophermart) GetAccount(ctx context.Context, user storage.User) (*storage.Account, error) {
	return g.container.Storage.LoadAccount(ctx, user)
}
