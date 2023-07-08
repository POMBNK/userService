package user

import (
	"context"
	"github.com/POMBNK/restAPI/pkg/logger"
)

type Service interface {
	Create(ctx context.Context, dto ToCreateUserDTO) (string, error)
	GetById(ctx context.Context, id string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, dto ToUpdateUserDTO) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	storage Storage
	logs    *logger.Logger
}

func (s service) Create(ctx context.Context, dto ToCreateUserDTO) (string, error) {
	user := CreateUserDto(dto)
	return s.storage.Create(ctx, user)
}

func (s service) GetById(ctx context.Context, id string) (User, error) {
	return s.storage.GetById(ctx, id)
}

func (s service) GetAll(ctx context.Context) ([]User, error) {
	return s.storage.GetAll(ctx)
}
func (s service) Update(ctx context.Context, dto ToUpdateUserDTO) error {
	user := UpdateUserDto(dto)
	return s.storage.Update(ctx, user)
}

func (s service) Delete(ctx context.Context, id string) error {
	return s.storage.Delete(ctx, id)
}

func New(storage Storage, logs *logger.Logger) Service {
	return &service{
		storage: storage,
		logs:    logs,
	}
}
