package gophermart

import (
	"context"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

// Register Registers a new user with given login and password
func (g Gophermart) Register(ctx context.Context, login string, rawPassword string) (*storage.User, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}
	if rawPassword == "" {
		return nil, ErrEmptyPassword
	}

	userID := utils.NewUUID6()
	user := storage.NewUser(userID, login, rawPassword)

	if err := g.container.Storage.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &user, nil
}
