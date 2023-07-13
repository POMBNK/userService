package db

import (
	"context"
	"fmt"
	"github.com/POMBNK/restAPI/internal/author"
	"github.com/POMBNK/restAPI/internal/pkg/apierror"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (d *postgresDB) Create(ctx context.Context, author author.Author) (string, error) {
	q := `INSERT INTO authors (name, surname) VALUES ($1,$2) RETURNING id`
	if err := d.client.QueryRow(ctx, q).Scan(&author.Id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("can not create user due error:%w", err)
	}
	return author.Id, nil
}

func (d *postgresDB) GetById(ctx context.Context, id string) (author.Author, error) {
	q := `SELECT id,name,surname FROM authors WHERE id = $1`
	var authorUnit author.Author
	err := d.client.QueryRow(ctx, q, id).Scan(&authorUnit.Id, authorUnit.Name, authorUnit.SurName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return author.Author{}, apierror.ErrNotFound
		}
		return author.Author{}, err
	}
	return authorUnit, nil
}

func (d *postgresDB) GetAll(ctx context.Context) ([]author.Author, error) {
	q := `SELECT id,name,surname FROM authors`
	rows, err := d.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	allAuthors := make([]author.Author, 0)
	for rows.Next() {
		var authorUnit author.Author
		err = rows.Scan(&authorUnit.Id, &authorUnit.Name, &authorUnit.SurName)
		if err != nil {
			return nil, err
		}
		allAuthors = append(allAuthors, authorUnit)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	//if rows.CommandTag().RowsAffected() != 1 {
	//	return nil, apierror.ErrNotFound
	//}
	return allAuthors, nil
}

func (d *postgresDB) Update(ctx context.Context, author author.Author) error {
	q := `UPDATE authors SET name=$1, surname=$2 WHERE id = $3 `
	//TODO:result from exec to logs.
	res, err := d.client.Exec(ctx, q, author.Name, author.SurName, author.Id)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}
	return nil
}

func (d *postgresDB) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM authors WHERE id = $1`
	res, err := d.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}
	return nil
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) author.Storage {
	return &postgresDB{
		client: client,
		logs:   logs}
}
