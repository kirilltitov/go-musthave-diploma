package storage

import "context"

type Storage interface {
	CreateUser(ctx context.Context, user User) error
}
