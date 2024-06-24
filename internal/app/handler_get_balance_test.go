package app

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestApplication_HandlerGetBalance(t *testing.T) {
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
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 401,
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
					result := storage.Account{
						UserID:           userID,
						CurrentBalance:   decimal.NewFromFloat(13.37),
						WithdrawnBalance: decimal.NewFromFloat(3.22),
					}
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadAccount(mock.Anything, mock.Anything).
						Return(&result, nil)
					return s
				}(),
			},
			want: want{
				code: 200,
				body: `{"current": 13.37, "withdrawn": 3.22}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Container.Storage = tt.input.storage
			r := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			if tt.input.cookie != nil {
				r.AddCookie(tt.input.cookie)
			}
			w := httptest.NewRecorder()

			a.HandlerGetBalance(w, r)

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
