package gophermart

import "errors"

var ErrEmptyLogin = errors.New("login string is empty")
var ErrEmptyPassword = errors.New("password string is empty")
var ErrAuthFailed = errors.New("wrong login or password")
var ErrInvalidOrderNumber = errors.New("invalid order number")
var ErrNotYourOrder = errors.New("not your order")
var ErrOrderAlreadyUploaded = errors.New("order already uploaded")
