package app

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-musthave-diploma/internal/accrual"
	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/container"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func TestApplication_HandlerCreateOrder(t *testing.T) {
	cfg := config.Config{}
	ctx := context.Background()

	accrualMock := accrual.NewMockAccrual(t)
	accrualMock.EXPECT().CalculateAmount(mock.Anything).Return(nil, accrual.ErrNoOrder)
	cnt := container.Container{Storage: nil, Accrual: accrualMock}

	a, err := New(ctx, cfg, &cnt)
	require.NoError(t, err)

	userID := utils.NewUUID6()
	user := storage.NewUser(userID, "frankstrino", "hesoyam")

	invalidOrderNumber := "1111"
	validOrderNumber := "79927398713"

	type input struct {
		cookie  *http.Cookie
		body    *string
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
				body:    nil,
				storage: storage.NewMockStorage(t),
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Negative (no body)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body:    nil,
				storage: storage.NewMockStorage(t),
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "Negative (invalid order number)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body:    &invalidOrderNumber,
				storage: storage.NewMockStorage(t),
			},
			want: want{
				code: 422,
			},
		},
		{
			name: "Negative (not your order)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body: &validOrderNumber,
				storage: func() storage.Storage {
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						LoadOrder(mock.Anything, mock.Anything).
						Return(&storage.Order{UserID: utils.NewUUID6()}, nil)
					return s
				}(),
			},
			want: want{
				code: 409,
			},
		},
		{
			name: "Positive",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body: &validOrderNumber,
				storage: func() storage.Storage {
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						LoadOrder(mock.Anything, mock.Anything).
						Return(nil, nil)
					s.
						EXPECT().
						CreateOrder(mock.Anything, mock.Anything).
						Return(errors.New("some error"))
					return s
				}(),
			},
			want: want{
				code: 500,
			},
		},
		{
			name: "Positive (already created)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body: &validOrderNumber,
				storage: func() storage.Storage {
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						LoadOrder(mock.Anything, mock.Anything).
						Return(&storage.Order{UserID: userID}, nil)
					return s
				}(),
			},
			want: want{
				code: 200,
			},
		},
		{
			name: "Positive",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				body: &validOrderNumber,
				storage: func() storage.Storage {
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						LoadOrder(mock.Anything, mock.Anything).
						Return(nil, nil)
					s.
						EXPECT().
						CreateOrder(mock.Anything, mock.Anything).
						Return(nil)
					return s
				}(),
			},
			want: want{
				code: 202,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Container.Storage = tt.input.storage
			var body []byte
			if tt.input.body != nil {
				body = []byte(*tt.input.body)
			}
			r := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewReader(body))
			if tt.input.cookie != nil {
				r.AddCookie(tt.input.cookie)
			}
			w := httptest.NewRecorder()

			a.HandlerCreateOrder(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
		})
	}
}
