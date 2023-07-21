package book

import (
	"context"
	"github.com/POMBNK/restAPI/pkg/logger"
)

type Service interface {
	Create(ctx context.Context, dto ToCreateBookDTO) (string, error)
}

type service struct {
	storage Storage
	logs    *logger.Logger
}

func (s *service) Create(ctx context.Context, dto ToCreateBookDTO) (string, error) {

	book := CreateBookDto(dto)
	id, err := s.storage.Create(ctx, book, book.AuthorID.Name, book.AuthorID.SurName)
	if err != nil {
		return "", err
	}
	return id, nil
}

func NewService(storage Storage, logs *logger.Logger) Service {
	return &service{
		storage: storage,
		logs:    logs,
	}
}
