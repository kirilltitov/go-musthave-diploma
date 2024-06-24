package gophermart

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

// GetOrders Returns orders for given user
func (g Gophermart) GetOrders(ctx context.Context, user storage.User) (*[]storage.Order, error) {
	return g.container.Storage.LoadOrders(ctx, user)
}
