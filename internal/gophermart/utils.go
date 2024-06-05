package gophermart

import (
	"strconv"

	"github.com/theplant/luhn"
)

func validateOrderNumber(order string) error {
	orderInt, err := strconv.Atoi(order)
	if err != nil || !luhn.Valid(orderInt) {
		return ErrInvalidOrder
	}

	return nil
}
