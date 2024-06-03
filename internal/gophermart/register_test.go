package gophermart

import (
	"context"
	"testing"

	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/container"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplication_HandlerRegister(t *testing.T) {
	cfg := config.New()
	ctx := context.Background()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	type input struct {
		login    string
		password string
		storage  storage.Storage
	}
	type want struct {
		err  error
		user *storage.User
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Negative (empty login)",
			input: input{
				password: "some",
			},
			want: want{
				err: ErrEmptyLogin,
			},
		},
		{
			name: "Negative (empty password)",
			input: input{
				login: `some`,
			},
			want: want{
				err: ErrEmptyPassword,
			},
		},
		{
			name: "Positive",
			input: input{
				login:    "frankstrino",
				password: "hesoyam",
				storage: func() storage.Storage {
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						CreateUser(mock.Anything, mock.Anything).
						Return(nil)
					return s
				}(),
			},
			want: want{
				user: &storage.User{
					Login:    "frankstrino",
					Password: "hesoyam",
				},
			},
		},
		{
			name: "Negative (duplicate)",
			input: input{
				login:    "frankstrino",
				password: "hesoyam",
				storage: func() storage.Storage {
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						CreateUser(mock.Anything, mock.Anything).
						Return(storage.ErrDuplicateFound)
					return s
				}(),
			},
			want: want{
				err: storage.ErrDuplicateFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.container.Storage = tt.input.storage

			user, err := g.Register(ctx, tt.input.login, tt.input.password)

			if tt.want.err != nil {
				assert.ErrorIs(t, tt.want.err, err)
			}
			if tt.want.user != nil {
				assert.NotNil(t, user)
				assert.Equal(t, tt.want.user.Login, user.Login)
			}
		})
	}
}
