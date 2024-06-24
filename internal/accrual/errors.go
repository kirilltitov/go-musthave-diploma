package accrual

import "errors"

var ErrNoOrder = errors.New("no such order")
var ErrInternalError = errors.New("internal accrual system error")

type ErrRateLimit int

func (e ErrRateLimit) Error() string {
	return "rate limit exceeded"
}
