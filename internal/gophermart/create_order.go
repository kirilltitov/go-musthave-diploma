package gophermart

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kirilltitov/go-musthave-diploma/internal/accrual"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

const (
	retries = 10
	timeout = 10
)

type errRetryWithTimeout struct {
	Timeout time.Duration
}

func (e errRetryWithTimeout) Error() string {
	return fmt.Sprintf("we should retry with timeout of %d", e.Timeout)
}

var errRetry = errors.New("retry")

var ctx = context.Background()

// CreateOrder Creates a new order in NEW state for given order number and then (in a non-blocking manner)
// queries the external accrual system for bonus amount
func (g Gophermart) CreateOrder(ctx context.Context, user storage.User, orderNumber string) error {
	if err := checkOrderNumber(orderNumber); err != nil {
		return err
	}

	existingOrder, err := g.container.Storage.LoadOrder(ctx, orderNumber)
	if err != nil {
		return err
	}
	if existingOrder != nil {
		if existingOrder.UserID != user.ID {
			return ErrNotYourOrder
		}
		return ErrOrderAlreadyUploaded
	}

	order := storage.Order{
		ID:          utils.NewUUID6(),
		OrderNumber: orderNumber,
		UserID:      user.ID,
		Status:      storage.StatusNew,
		CreatedAt:   time.Now(),
	}
	if err := g.container.Storage.CreateOrder(ctx, order); err != nil {
		return err
	}

	g.acquireBalanceForOrder(user, order)

	return nil
}

func (g Gophermart) acquireBalanceForOrder(user storage.User, order storage.Order) {
	logger := utils.Log

	go func() {
		for i := 0; i < retries; i++ {

			result, err := g.container.Accrual.CalculateAmount(order)
			if err != nil {
				if errors.Is(err, accrual.ErrNoOrder) {
					logger.Errorf("Order %s not found in accrual system", order.OrderNumber)
					return
				} else if e, ok := err.(accrual.ErrRateLimit); ok {
					logger.Errorf("Rate limit error, retrying in %d seconds", e)
					time.Sleep(time.Duration(e) * time.Second)
					continue
				} else if errors.Is(err, accrual.ErrInternalError) {
					logger.Errorf("Internal error from accrual system: %v", err)
					continue
				} else {
					logger.Errorf("Unexpected error from accrual system: %v", err)
					continue
				}
			}

			if err := g.processAndApplyAccrualResponseToOrder(ctx, order, user, result); err != nil {
				var e errRetryWithTimeout
				if errors.As(err, &e) {
					time.Sleep(e.Timeout)
					continue
				} else if errors.Is(err, errRetry) {
					continue
				} else {
					logger.Errorf("Unknown error while applying accrual response to order: %+v", err)
					return
				}
			}

			return
		}
	}()
}

func (g Gophermart) processAndApplyAccrualResponseToOrder(ctx context.Context, order storage.Order, user storage.User, result *accrual.CalculationResult) error {
	logger := utils.Log

	switch result.Status {
	case accrual.StatusRegistered:
		logger.Infof("Order is registered, but not processed, retrying in %d seconds", timeout)
		time.Sleep(timeout * time.Second)
		return errRetry
	case accrual.StatusProcessing:
		logger.Infof("Order is processing, retrying in %d seconds", timeout)

		newStatus := storage.StatusProcessing
		if err := g.container.Storage.UpdateOrderStatus(
			ctx,
			order,
			newStatus,
			[]storage.OrderStatus{storage.StatusNew, storage.StatusProcessing},
		); err != nil {
			logger.Errorf(
				"Error occurred while trying to update order status from %s to %s: %s, exiting",
				order.Status, newStatus, err,
			)
			return nil
		}

		time.Sleep(timeout * time.Second)
		return errRetry
	case accrual.StatusInvalid:
		logger.Infof("Order is invalid, exiting")

		newStatus := storage.StatusInvalid
		if err := g.container.Storage.UpdateOrderStatus(
			ctx,
			order,
			newStatus,
			[]storage.OrderStatus{},
		); err != nil {
			logger.Errorf(
				"Error occurred while trying to update order status to %s: %s",
				newStatus, err,
			)
			return nil
		}

		return nil
	case accrual.StatusProcessed:
		logger.Infof("Order has been processed, saving")

		if result.AccrualRaw == nil {
			logger.Errorf(
				"Could not apply acrrual amount to account because accrual response does not contain amount: %v",
				result,
			)
			return nil
		}

		if err := g.container.Storage.ApplyProcessedOrder(
			ctx,
			user,
			order,
			*result.Accrual(),
		); err != nil {
			logger.Errorf(
				"Error occurred while trying to apply accrual amount to account: %s",
				err,
			)
			return nil
		}

		return nil
	}

	return nil
}
