package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/POMBNK/restAPI/internal/book"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (p *postgresDB) Create(ctx context.Context, book book.Book, name, surname string) (id string, err error) {
	q := `SELECT DISTINCT id FROM authors WHERE name=$1 AND surname = $2`
	err = p.client.QueryRow(ctx, q, name, surname).Scan(&book.AuthorID.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			book.AuthorID.Name = name
			book.AuthorID.SurName = surname
		} else {
			return "", err
		}
	}
	//TODO:
	// BEGIN TX
	// CREATE AUTHOR QUERY WITH ID,name,surname
	// CREATE book with this authorID
	tx, err := p.client.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed begin Tx due error:%w", err)
	}
	q = `INSERT INTO authors (name, surname) VALUES ($1,$2) RETURNING id`
	err = tx.QueryRow(ctx, q, book.AuthorID.Name, book.AuthorID.SurName).Scan(&book.AuthorID.Id)
	if err != nil {
		tx.Rollback(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("can not create user due error:%w", err)
	}
	q = `INSERT INTO books (name, author_id)  VALUES ($1,$2) RETURNING id`
	err = tx.QueryRow(ctx, q, book.Name, book.AuthorID.Id).Scan(&book.Id)
	if err != nil {
		tx.Rollback(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("can not create user due error:%w", err)
	}

	return id, tx.Commit(ctx)
}

func (p *postgresDB) GetByID(ctx context.Context, id string) (book.Book, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postgresDB) GetByAuthor(ctx context.Context, author string) {
	//TODO implement me
	panic("implement me")
}

func (p *postgresDB) GetByName(ctx context.Context, userName string) {
	//TODO implement me
	panic("implement me")
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) book.Storage {
	return &postgresDB{
		client: client,
		logs:   logs,
	}
}
