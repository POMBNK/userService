package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/POMBNK/restAPI/internal/auth"
	"github.com/POMBNK/restAPI/internal/pkg/apierror"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/jackc/pgx/v5"
	"time"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (d *postgresDB) Create(ctx context.Context, user auth.User) (string, error) {
	q := `INSERT INTO users (username, password_hash, email)  VALUES ($1,$2,$3) RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := d.client.QueryRow(ctx, q, user.UserName,
		user.PasswordHash, user.Email).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("can not create user due error:%w", err)
	}

	return user.ID, nil
}

func (d *postgresDB) GetById(ctx context.Context, id string) (auth.User, error) {
	var userUnit auth.User
	q := `SELECT username,password_hash,email FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := d.client.QueryRow(ctx, q, id).Scan(&userUnit.UserName, &userUnit.PasswordHash, &userUnit.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.User{}, apierror.ErrNotFound
		}
		return auth.User{}, err
	}
	return userUnit, nil
}

func (d *postgresDB) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	var userUnit auth.User
	q := `SELECT id,username,password_hash,email FROM users WHERE email = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := d.client.QueryRow(ctx, q, email).Scan(&userUnit.ID, &userUnit.UserName, &userUnit.PasswordHash, &userUnit.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.User{}, apierror.ErrNotFound
		}
		return auth.User{}, err
	}
	return userUnit, nil
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) auth.Storage {
	return &postgresDB{
		client: client,
		logs:   logs}
}
