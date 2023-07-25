package book

import (
	"context"
	"github.com/POMBNK/restAPI/pkg/logger"
)

type Service interface {
	Create(ctx context.Context, dto ToCreateBookDTO) (string, error)
	GetByID(ctx context.Context, id string) (Book, error)
	GetByAuthor(ctx context.Context, authorID string) ([]Book, error)
	GetByName(ctx context.Context, name string) ([]Book, error)
}

type service struct {
	storage Storage
	logs    *logger.Logger
}

func (s *service) Create(ctx context.Context, dto ToCreateBookDTO) (string, error) {

	book := CreateBookDto(dto)
	// TODO: Check has book already exist before
	id, err := s.storage.Create(ctx, book)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *service) GetByID(ctx context.Context, id string) (Book, error) {
	//TODO: Fill all fields in response
	book, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return Book{}, err
	}
	return book, nil
}

func (s *service) GetByAuthor(ctx context.Context, authorID string) ([]Book, error) {
	//TODO: Fill all fields in response
	books, err := s.storage.GetByAuthor(ctx, authorID)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *service) GetByName(ctx context.Context, name string) ([]Book, error) {
	//TODO: Fill all fields in response
	books, err := s.storage.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func NewService(storage Storage, logs *logger.Logger) Service {
	return &service{
		storage: storage,
		logs:    logs,
	}
}
