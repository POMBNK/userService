package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/POMBNK/restAPI/internal/book"
	"github.com/POMBNK/restAPI/internal/pkg/apierror"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (p *postgresDB) Create(ctx context.Context, book book.Book) (id string, err error) {
	q := `SELECT DISTINCT id FROM authors WHERE name=$1 AND surname = $2`
	err = p.client.QueryRow(ctx, q, book.AuthorID.Name, book.AuthorID.SurName).Scan(&book.AuthorID.Id)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
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
	var bookUnit book.Book
	q := `SELECT name, author_id FROM books WHERE id = $1`
	err := p.client.QueryRow(ctx, q, id).Scan(&bookUnit.Name, &bookUnit.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return book.Book{}, apierror.ErrNotFound
		}
		return book.Book{}, err
	}

	return bookUnit, nil
}

func (p *postgresDB) GetByAuthor(ctx context.Context, authorID string) ([]book.Book, error) {
	q := `SELECT name FROM books WHERE author_id = $1`
	booksRow, err := p.client.Query(ctx, q, authorID)
	if err != nil {
		return nil, err
	}
	books := make([]book.Book, 0)
	for booksRow.Next() {
		var bookUnit book.Book
		err = booksRow.Scan(&bookUnit.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, bookUnit)
	}
	if err = booksRow.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (p *postgresDB) GetByName(ctx context.Context, bookName string) ([]book.Book, error) {
	q := `SELECT id, author_id FROM books WHERE name = $1`
	booksRow, err := p.client.Query(ctx, q, bookName)
	if err != nil {
		return nil, err
	}
	books := make([]book.Book, 0)
	for booksRow.Next() {
		var bookUnit book.Book
		err = booksRow.Scan(&bookUnit.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, bookUnit)
	}
	if err = booksRow.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) book.Storage {
	return &postgresDB{
		client: client,
		logs:   logs,
	}
}
