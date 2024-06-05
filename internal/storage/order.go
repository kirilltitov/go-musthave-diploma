package storage

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

type Order struct {
	ID          uuid.UUID       `db:"id"`
	OrderNumber string          `db:"order_number"`
	UserID      uuid.UUID       `db:"created_at"`
	Status      string          `db:"status"`
	Amount      decimal.Decimal `db:"amount"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   *time.Time      `db:"updated_at"`
}

func (s PgSQL) LoadOrders(ctx context.Context, user User) (*[]Order, error) {
	var rows []Order

	err := pgxscan.Select(ctx, s.Conn, &rows, `select * from public.orders where user_id = $1 order by created_at asc`, user.ID)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}
