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
	"time"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (d *postgresDB) Create(ctx context.Context, author author.Author) (string, error) {
	q := `INSERT INTO authors (name, surname) VALUES ($1,$2) RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := d.client.QueryRow(ctx, q).Scan(&author.Id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("can not create user due error:%w", err)
	}
	d.logs.Debug("Author created")
	d.logs.Tracef("id of created user: %s \n", author.Id)

	return author.Id, nil
}

func (d *postgresDB) GetById(ctx context.Context, id string) (author.Author, error) {
	var authorUnit author.Author

	q := `SELECT id,name,surname FROM authors WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err := d.client.QueryRow(ctx, q, id).Scan(&authorUnit.Id, &authorUnit.Name, &authorUnit.SurName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return author.Author{}, apierror.ErrNotFound
		}
		return author.Author{}, err
	}

	d.logs.Debugf("Author GetByID success")

	return authorUnit, nil
}

func (d *postgresDB) GetAll(ctx context.Context) ([]author.Author, error) {
	q := `SELECT id,name,surname FROM authors`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
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

	d.logs.Debugf("Author GetAll success")

	return allAuthors, nil
}

func (d *postgresDB) Update(ctx context.Context, author author.Author) error {
	q := `UPDATE authors SET name=$1, surname=$2 WHERE id = $3 `
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	//TODO:result from exec to logs.
	res, err := d.client.Exec(ctx, q, author.Name, author.SurName, author.Id)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}

	d.logs.Tracef("Matched and updated %v documents.\n", res.RowsAffected())

	return nil
}

func (d *postgresDB) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM authors WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	res, err := d.client.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}

	d.logs.Tracef("Deleted %v documents.\n", res.RowsAffected())

	return nil
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) author.Storage {
	return &postgresDB{
		client: client,
		logs:   logs}
}
