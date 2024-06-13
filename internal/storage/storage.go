package storage

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockery
type Storage interface {
	CreateUser(ctx context.Context, user User) error
	LoadUser(ctx context.Context, login string) (*User, error)
	CreateOrder(ctx context.Context, order Order) error
	LoadOrder(ctx context.Context, orderNumber string) (*Order, error)
	LoadOrders(ctx context.Context, user User) (*[]Order, error)
	UpdateOrderStatus(ctx context.Context, order Order, newStatus OrderStatus, allowedOldStatuses []OrderStatus) error
	ApplyProcessedOrder(ctx context.Context, user User, order Order, amount decimal.Decimal) error
	LoadAccount(ctx context.Context, user User) (*Account, error)
	LoadWithdrawals(ctx context.Context, user User) (*[]Withdrawal, error)
	WithdrawBalanceFromAccount(ctx context.Context, user User, amount decimal.Decimal, order string) error
}
