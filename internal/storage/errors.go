package storage

import "errors"

var ErrNotFound = errors.New("key was not found")
var ErrDuplicateFound = errors.New("user with this login already exists")
var ErrDeleted = errors.New("URL has been deleted")
var ErrInsufficientBalance = errors.New("insufficient balance")
