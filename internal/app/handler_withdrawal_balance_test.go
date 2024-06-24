package app

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/container"
	mockStorage "github.com/kirilltitov/go-musthave-diploma/internal/mocks/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func TestApplication_HandlerWithdrawBalance(t *testing.T) {
	cfg := config.Config{}
	cfg.JWTCookieName = "access_token"
	cfg.JWTSecret = "hesoyam"
	cfg.JWTTimeToLive = 86400
	ctx := context.Background()
	cnt := container.Container{Storage: nil}

	a, err := New(ctx, cfg, &cnt)
	require.NoError(t, err)

	userID := utils.NewUUID6()
	user := storage.NewUser(userID, "frankstrino", "hesoyam")

	type input struct {
		cookie  *http.Cookie
		body    string
		storage storage.Storage
	}
	type want struct {
		code int
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Negative (unauthorized)",
			input: input{
				cookie:  nil,
				body:    `{}`,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Negative (invalid order 1)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body:    `{"order":"lul","sum":0}`,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 422,
			},
		},
		{
			name: "Negative (invalid order 2)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body:    `{"order":"1111","sum":0}`,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 422,
			},
		},
		{
			name: "Negative (insufficient balance)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body: `{"order":"79927398713","sum":1337}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						WithdrawBalanceFromAccount(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(storage.ErrInsufficientBalance)
					return s
				}(),
			},
			want: want{
				code: 402,
			},
		},
		{
			name: "Positive",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body: `{"order":"79927398713","sum":322}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						WithdrawBalanceFromAccount(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
						Return(nil)
					return s
				}(),
			},
			want: want{
				code: 200,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Container.Storage = tt.input.storage
			r := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", bytes.NewReader([]byte(tt.input.body)))
			if tt.input.cookie != nil {
				r.AddCookie(tt.input.cookie)
			}
			w := httptest.NewRecorder()

			a.HandlerWithdrawBalance(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
		})
	}
}
