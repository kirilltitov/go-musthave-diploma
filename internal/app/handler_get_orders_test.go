package app

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kinbiko/jsonassert"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/go-musthave-diploma/internal/config"
	"github.com/kirilltitov/go-musthave-diploma/internal/container"
	mockStorage "github.com/kirilltitov/go-musthave-diploma/internal/mocks/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func TestApplication_HandlerGetOrders(t *testing.T) {
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
		storage storage.Storage
	}
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Negative (unathorized)",
			input: input{
				cookie:  nil,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Positive (no results)",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				storage: func() storage.Storage {
					orders := make([]storage.Order, 0)
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadOrders(mock.Anything, mock.Anything).
						Return(&orders, nil)
					return s
				}(),
			},
			want: want{
				code: 204,
			},
		},
		{
			name: "Positive",
			input: input{
				cookie: func() *http.Cookie {
					cookie, _ := a.createAuthCookie(user)
					return cookie
				}(),
				storage: func() storage.Storage {
					orders := []storage.Order{
						{
							ID:          utils.NewUUID6(),
							OrderNumber: "123",
							UserID:      userID,
							Status:      storage.StatusNew,
							Amount: func() *decimal.Decimal {
								r := decimal.NewFromFloat(13.37)
								return &r
							}(),
							CreatedAt: time.Now(),
							UpdatedAt: nil,
						},
						{
							ID:          utils.NewUUID6(),
							OrderNumber: "456",
							UserID:      userID,
							Status:      storage.StatusProcessed,
							Amount: func() *decimal.Decimal {
								r := decimal.NewFromFloat(3.22)
								return &r
							}(),
							CreatedAt: time.Now(),
							UpdatedAt: nil,
						},
					}
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadOrders(mock.Anything, mock.Anything).
						Return(&orders, nil)
					return s
				}(),
			},
			want: want{
				code: 200,
				body: `[
					{
						"number": "123",
						"status": "NEW",
						"accrual": 13.37,
						"uploaded_at": "<<PRESENCE>>"
					},
					{
						"number": "456",
						"status": "PROCESSED",
						"accrual": 3.22,
						"uploaded_at": "<<PRESENCE>>"
					}
				]`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Container.Storage = tt.input.storage
			r := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			if tt.input.cookie != nil {
				r.AddCookie(tt.input.cookie)
			}
			w := httptest.NewRecorder()

			a.HandlerGetOrders(w, r)

			result := w.Result()
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			require.Equal(t, tt.want.code, result.StatusCode)
			ja := jsonassert.New(t)
			ja.Assertf(string(resultBody), tt.want.body)
		})
	}
}
