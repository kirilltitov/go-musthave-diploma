package storage

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUser(ID uuid.UUID, login string, rawPassword string) User {
	user := User{
		ID:        ID,
		Login:     login,
		CreatedAt: time.Time{},
	}

	h := sha1.New()
	io.WriteString(h, rawPassword)
	io.WriteString(h, strconv.FormatInt(user.CreatedAt.Unix(), 10))

	user.Password = hex.EncodeToString(h.Sum(nil))

	return user
}

func LoadUser(ctx context.Context, querier pgxscan.Querier, login string) (*User, error) {
	var row User

	if err := pgxscan.Get(ctx, querier, &row, `select * from public.user where login = $1`, login); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &row, nil
}

func insertUser(ctx context.Context, querier pgxscan.Querier, user User) error {
	query := `insert into public.user (id, login, password, created_at) values ($1, $2, $3, $4)`
	_, err := querier.Query(ctx, query, user.ID, user.Login, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s PgSQL) CreateUser(ctx context.Context, user User) error {
	_, err := WithTransaction(ctx, s, func(tx pgx.Tx) (*any, error) {
		existingUser, err := LoadUser(ctx, tx, user.Login)
		if err != nil {
			return nil, err
		}
		if existingUser != nil {
			return nil, ErrDuplicateFound
		}

		if err := insertUser(ctx, tx, user); err != nil {
			return nil, err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}