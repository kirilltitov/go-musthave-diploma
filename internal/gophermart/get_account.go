package gophermart

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

// GetAccount Returns account info for given user
func (g Gophermart) GetAccount(ctx context.Context, user storage.User) (*storage.Account, error) {
	return g.container.Storage.LoadAccount(ctx, user)
}
