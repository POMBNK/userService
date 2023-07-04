package user

import (
	"context"
	"github.com/POMBNK/restAPI/pkg/logger"
)

type Service struct {
	storage Storage
	logs    *logger.Logger
}

func (s *Service) Create(ctx context.Context, dto UserDTO) (User, error) {
	//s.storage.Create(ctx,dto)
	return User{}, nil
}
