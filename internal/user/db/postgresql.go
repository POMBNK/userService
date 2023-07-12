package db

import (
	"context"
	"github.com/POMBNK/restAPI/internal/user"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/logger"
)

type postgresDB struct {
	client postgresql.Client
	logs   *logger.Logger
}

func (p postgresDB) Create(ctx context.Context, user user.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) GetById(ctx context.Context, id string) (user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) GetAll(ctx context.Context) ([]user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) Update(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (p postgresDB) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewPostgresDB(client postgresql.Client, logs *logger.Logger) user.Storage {
	return &postgresDB{
		client: client,
		logs:   logs}
}
