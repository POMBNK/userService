package db

import (
	"context"
	"github.com/POMBNK/restAPI/internal/user"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDB struct {
	pool *pgxpool.Pool
	logs *logger.Logger
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

func NewPostgresDB(pool *pgxpool.Pool) user.Storage {
	return &postgresDB{pool: pool}
}
