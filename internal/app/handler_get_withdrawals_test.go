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
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func TestApplication_HandlerGetWithdrawals(t *testing.T) {
	cfg := config.Config{}
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
				storage: storage.NewMockStorage(t),
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
					withdrawals := make([]storage.Withdrawal, 0)
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						LoadWithdrawals(mock.Anything, mock.Anything).
						Return(&withdrawals, nil)
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
					withdrawals := []storage.Withdrawal{
						{
							ID:          utils.NewUUID6(),
							UserID:      userID,
							OrderNumber: "123",
							Amount:      decimal.NewFromFloat(13.37),
							CreatedAt:   time.Now(),
						},
						{
							ID:          utils.NewUUID6(),
							UserID:      userID,
							OrderNumber: "456",
							Amount:      decimal.NewFromFloat(3.22),
							CreatedAt:   time.Now(),
						},
					}
					s := storage.NewMockStorage(t)
					s.
						EXPECT().
						LoadWithdrawals(mock.Anything, mock.Anything).
						Return(&withdrawals, nil)
					return s
				}(),
			},
			want: want{
				code: 200,
				body: `[
					{"order": "123","sum": 13.37,"processed_at": "<<PRESENCE>>"},
					{"order": "456","sum": 3.22,"processed_at": "<<PRESENCE>>"}
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

			a.HandlerGetWithdrawals(w, r)

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
