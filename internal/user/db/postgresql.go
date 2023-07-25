package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/POMBNK/restAPI/internal/pkg/apierror"
	"github.com/POMBNK/restAPI/internal/user"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/jackc/pgx/v5"
	"time"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (d *postgresDB) Create(ctx context.Context, user user.User) (string, error) {
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

func (d *postgresDB) GetById(ctx context.Context, id string) (user.User, error) {
	var userUnit user.User
	q := `SELECT username,password_hash,email FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := d.client.QueryRow(ctx, q, id).Scan(&userUnit.UserName, &userUnit.PasswordHash, &userUnit.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user.User{}, apierror.ErrNotFound
		}
		return user.User{}, err
	}
	return userUnit, nil
}

func (d *postgresDB) GetAll(ctx context.Context) ([]user.User, error) {
	q := `SELECT id,username,password_hash,email FROM users`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	rows, err := d.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	allUsers := make([]user.User, 0)
	for rows.Next() {
		var userUnit user.User
		err = rows.Scan(&userUnit.ID, &userUnit.UserName, &userUnit.PasswordHash, &userUnit.Email)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, userUnit)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return allUsers, nil
}

func (d *postgresDB) Update(ctx context.Context, user user.User) error {
	q := `UPDATE users SET username=$1, password_hash=$2,email=$3 WHERE id = $4 `
	//TODO:result from exec to logs.
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	res, err := d.client.Exec(ctx, q, user.UserName, user.PasswordHash, user.Email, user.ID)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}
	return nil
}

func (d *postgresDB) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	res, err := d.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}
	return nil
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) user.Storage {
	return &postgresDB{
		client: client,
		logs:   logs}
}
