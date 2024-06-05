package storage

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockery
type Storage interface {
	CreateUser(ctx context.Context, user User) error
	LoadUser(ctx context.Context, login string) (*User, error)
	LoadOrders(ctx context.Context, user User) (*[]Order, error)
	LoadAccount(ctx context.Context, user User) (*Account, error)
	LoadWithdrawals(ctx context.Context, user User) (*[]Withdrawal, error)
	WithdrawBalanceFromAccount(ctx context.Context, user User, amount decimal.Decimal, order string) error
}
