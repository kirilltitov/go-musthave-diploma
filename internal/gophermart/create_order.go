package gophermart

import (
	"context"
	"errors"
	"time"

	"github.com/kirilltitov/go-musthave-diploma/internal/accrual"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func (g Gophermart) CreateOrder(ctx context.Context, user storage.User, orderNumber string) error {
	if err := validateOrderNumber(orderNumber); err != nil {
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

	g.acquireBalanceForOrder(ctx, user, order)

	return nil
}

func (g Gophermart) acquireBalanceForOrder(ctx context.Context, user storage.User, order storage.Order) {
	logger := utils.Log
	retries := 10
	var timeout time.Duration = 10

	go func() {
		for i := 0; i < retries; i++ {
			result, err := g.container.Accrual.CalculateAmount(order)
			if err != nil {
				if errors.Is(err, accrual.ErrNoOrder) {
					logger.Infof("Order %s not found in accrual system", order.OrderNumber)
					return
				} else if e, ok := err.(accrual.ErrRateLimit); ok {
					logger.Infof("Rate limit error, retrying in %d seconds", e)
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

			switch result.Status {
			case accrual.StatusRegistered:
				logger.Infof("Order is registered, but not processed, retrying in %d seconds", timeout)
				time.Sleep(timeout * time.Second)
				continue
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
					return
				}

				time.Sleep(timeout * time.Second)
				continue
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
					return
				}

				return
			case accrual.StatusProcessed:
				logger.Infof("Order has been processed, saving")

				if result.Accrual == nil {
					logger.Errorf(
						"Could not apply acrrual amount to account because accrual response does not contain amount: %v",
						result,
					)
					return
				}

				if err := g.container.Storage.ApplyProcessedOrder(
					ctx,
					user,
					order,
					*result.Accrual,
				); err != nil {
					logger.Errorf(
						"Error occurred while trying to apply accrual amount to account: %s",
						err,
					)
					return
				}

				return
			}
		}
	}()
}
