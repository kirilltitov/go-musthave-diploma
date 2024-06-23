package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
	"github.com/shopspring/decimal"
)

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID          uuid.UUID        `db:"id"`
	OrderNumber string           `db:"order_number"`
	UserID      uuid.UUID        `db:"user_id"`
	Status      OrderStatus      `db:"status"`
	Amount      *decimal.Decimal `db:"amount"`
	CreatedAt   time.Time        `db:"created_at"`
	UpdatedAt   *time.Time       `db:"updated_at"`
}

func (s PgSQL) LoadOrders(ctx context.Context, user User) (*[]Order, error) {
	var rows []Order

	err := pgxscan.Select(ctx, s.Conn, &rows, `select * from public.order where user_id = $1 order by created_at asc`, user.ID)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

func (s PgSQL) LoadOrder(ctx context.Context, orderNumber string) (*Order, error) {
	return WithTransaction(ctx, s, func(tx pgx.Tx) (*Order, error) {
		return loadOrder(ctx, tx, orderNumber, false, false)
	})
}

func (s PgSQL) CreateOrder(ctx context.Context, order Order) error {
	_, err := s.Conn.Query(
		ctx,
		`insert into public.order (id, order_number, user_id, status, created_at) values ($1, $2, $3, $4, $5)`,
		order.ID, order.OrderNumber, order.UserID, order.Status, order.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s PgSQL) UpdateOrderStatus(ctx context.Context, order Order, newStatus OrderStatus, allowedOldStatuses []OrderStatus) error {
	return WithVoidTransaction(ctx, s, func(tx pgx.Tx) error {
		if err := updateOrderStatus(ctx, tx, order, newStatus, allowedOldStatuses); err != nil {
			return err
		}

		return tx.Commit(ctx)
	})
}

func loadOrder(ctx context.Context, tx pgx.Tx, orderNumber string, forUpdate bool, skipLocked bool) (*Order, error) {
	var result Order

	var suffixes []string
	if forUpdate {
		suffixes = append(suffixes, "for update")
	}
	if skipLocked {
		suffixes = append(suffixes, "skip locked")
	}
	query := fmt.Sprintf(`select * from public.order where order_number = $1 %s`, strings.Join(suffixes, " "))
	if err := pgxscan.Get(ctx, tx, &result, query, orderNumber); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &result, nil
}

func updateOrderStatus(ctx context.Context, tx pgx.Tx, order Order, newStatus OrderStatus, allowedOldStatuses []OrderStatus) error {
	orderFromDB, err := loadOrder(ctx, tx, order.OrderNumber, true, true)
	if err != nil {
		return err
	}
	if orderFromDB == nil {
		return ErrNotFound
	}
	if len(allowedOldStatuses) > 0 && !utils.InArray(allowedOldStatuses, orderFromDB.Status) {
		utils.Log.Errorf(
			"Unexpected status '%s' for order %s (expected %v), exiting",
			orderFromDB.Status, orderFromDB.OrderNumber, allowedOldStatuses,
		)
		return ErrWrongStatus(orderFromDB.Status)
	}

	res, err := tx.Exec(
		ctx,
		`update public.order set status = $1, updated_at = NOW() where ID = $2`,
		newStatus, order.ID,
	)
	if err != nil {
		return err
	}

	utils.Log.Debugf("Order status update result: %+v", res)

	return nil
}

func updateOrderAmount(ctx context.Context, tx pgx.Tx, order Order) error {
	res, err := tx.Exec(
		ctx,
		`update public.order set amount = $1, updated_at = NOW() where ID = $2`,
		order.Amount.InexactFloat64(), order.ID,
	)
	if err != nil {
		return err
	}

	utils.Log.Debugf("Order amount update result: %+v", res)

	return nil
}
