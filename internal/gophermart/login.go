package gophermart

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

func (g Gophermart) Login(ctx context.Context, login string, password string) (*storage.User, error) {
	user, err := g.container.Storage.LoadUser(ctx, login)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsValidPassword(password) {
		return nil, ErrAuthFailed
	}

	return user, nil
}
