package author

import (
	"context"
	"github.com/POMBNK/restAPI/pkg/logger"
)

type Service interface {
	Create(ctx context.Context, dto ToCreateAuthorDTO) (string, error)
	GetById(ctx context.Context, id string) (Author, error)
	GetAll(ctx context.Context) ([]Author, error)
	Update(ctx context.Context, dto ToUpdateAuthorDTO) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	storage Storage
	logs    *logger.Logger
}

func (s *service) Create(ctx context.Context, dto ToCreateAuthorDTO) (string, error) {
	author := CreateAuthorDto(dto)

	uuid, err := s.storage.Create(ctx, author)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (s *service) GetById(ctx context.Context, id string) (Author, error) {
	author, err := s.storage.GetById(ctx, id)
	if err != nil {
		return Author{}, err
	}

	return author, nil
}

func (s *service) GetAll(ctx context.Context) ([]Author, error) {
	authors, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return authors, nil
}

func (s *service) Update(ctx context.Context, dto ToUpdateAuthorDTO) error {

	origAuthor, err := s.storage.GetById(ctx, dto.ID)
	if err != nil {
		return err
	}
	if dto.Name == "" {
		dto.Name = origAuthor.Name
	}
	if dto.SurName == "" {
		dto.SurName = origAuthor.SurName
	}
	author := UpdateAuthorDto(dto)
	if err = s.storage.Update(ctx, author); err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.storage.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func newService(storage Storage, logs *logger.Logger) Service {
	return &service{storage: storage, logs: logs}
}
