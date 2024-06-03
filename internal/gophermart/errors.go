package gophermart

import "errors"

var ErrEmptyLogin = errors.New("login string is empty")
var ErrEmptyPassword = errors.New("password string is empty")
