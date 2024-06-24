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
	"github.com/kirilltitov/go-musthave-diploma/internal/gophermart"
	mockStorage "github.com/kirilltitov/go-musthave-diploma/internal/mocks/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"
)

func TestApplication_HandlerLogin(t *testing.T) {
	cfg := config.New()
	ctx := context.Background()
	cnt := container.Container{Storage: nil}

	a, err := New(ctx, cfg, &cnt)
	require.NoError(t, err)

	type input struct {
		body    string
		storage storage.Storage
	}
	type want struct {
		code      int
		cookieSet bool
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Negative (invalid request 1)",
			input: input{
				body:    `invalid`,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "Negative (invalid request 2)",
			input: input{
				body: `{}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadUser(mock.Anything, mock.Anything).
						Return(nil, gophermart.ErrAuthFailed)
					return s
				}(),
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Negative (invalid request 3)",
			input: input{
				body: `{"login":"frankstrino","passworddd":"hesoyam"}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadUser(mock.Anything, mock.Anything).
						Return(nil, gophermart.ErrAuthFailed)
					return s
				}(),
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Positive",
			input: func() input {
				userID := utils.NewUUID6()
				user := storage.NewUser(userID, "frankstrino", "hesoyam")

				return input{
					body: `{"login":"frankstrino","password":"hesoyam"}`,
					storage: func() storage.Storage {
						s := mockStorage.NewMockStorage(t)
						s.
							EXPECT().
							LoadUser(mock.Anything, mock.Anything).
							Return(&user, nil)
						return s
					}(),
				}
			}(),
			want: want{
				code:      200,
				cookieSet: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Container.Storage = tt.input.storage
			r := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader([]byte(tt.input.body)))
			w := httptest.NewRecorder()

			a.HandlerLogin(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			if tt.want.cookieSet {
				assert.NotEmpty(t, result.Cookies())
			} else {
				assert.Empty(t, result.Cookies())
			}
		})
	}
}
