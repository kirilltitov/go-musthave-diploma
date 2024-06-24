package storage

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("key was not found")
var ErrDuplicateFound = errors.New("user with this login already exists")
var ErrInsufficientBalance = errors.New("insufficient balance")

type ErrWrongStatus string

func (e ErrWrongStatus) Error() string {
	return fmt.Sprintf("wrong order status: %s", string(e))
}
