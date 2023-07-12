package db

import (
	"context"
	"github.com/POMBNK/restAPI/internal/author"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/logger"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (p postgresDB) Create(ctx context.Context, user author.Author) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) GetById(ctx context.Context, id string) (author.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) GetAll(ctx context.Context) ([]author.Author, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) Update(ctx context.Context, user author.Author) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) author.Storage {
	return &postgresDB{
		client: client,
		logs:   logs}
}
