package app

import (
	"net/http"
	"time"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
	"github.com/kirilltitov/go-musthave-diploma/internal/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
}

type UserID uuid.UUID

const (
	JWTCookieName = "access_token"
	JWTSecret     = "hesoyam"
	JWTTimeToLive = 86400
)

func (a *Application) createAuthCookie(user storage.User) (*http.Cookie, error) {
	exp := time.Now().Add(time.Second * JWTTimeToLive)
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   user.ID.String(),
				ExpiresAt: jwt.NewNumericDate(exp),
			},
		},
	)

	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:    JWTCookieName,
		Value:   tokenString,
		Expires: exp,
	}

	return &cookie, nil
}

func (a *Application) authenticate(r *http.Request) (*storage.User, error) {
	cookie, err := r.Cookie(JWTCookieName)
	if err != nil {
		return nil, err
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	if claims.Subject == "" {
		return nil, nil
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	utils.Log.Infof("Authenticated user %s by JWT cookie", userID.String())

	return &storage.User{ID: userID}, nil
}
