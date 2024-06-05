package gophermart

import "errors"

var ErrEmptyLogin = errors.New("login string is empty")
var ErrEmptyPassword = errors.New("password string is empty")
var ErrAuthFailed = errors.New("wrong login or password")
var ErrInvalidOrder = errors.New("invalid order number")
