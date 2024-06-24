package gophermart

import (
	"strconv"

	"github.com/theplant/luhn"
)

func checkOrderNumber(order string) error {
	orderInt, err := strconv.ParseInt(order, 10, 64)
	if err != nil || !luhn.Valid(int(orderInt)) {
		return ErrInvalidOrderNumber
	}

	return nil
}
