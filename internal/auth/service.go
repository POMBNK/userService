package auth

import (
	"context"
	"fmt"
	"github.com/POMBNK/restAPI/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUP(ctx context.Context, dto ToCreateUserDTO) (string, error)
	SignIN(ctx context.Context, id string) (User, error)
}

type service struct {
	storage Storage
	logs    *logger.Logger
}

func (s *service) SignUP(ctx context.Context, dto ToCreateUserDTO) (string, error) {
	user := CreateUserDto(dto)

	pswrd, err := hashPassword(dto.Password)
	if err != nil {
		return "", fmt.Errorf("failled due hashing password error:%w", err)
	}
	user.PasswordHash = pswrd
	uuid, err := s.storage.Create(ctx, user)
	if err != nil {
		return uuid, err
	}
	//TODO: create pair of JWT token
	return uuid, nil
}

func (s *service) SignIN(ctx context.Context, id string) (User, error) {
	//TODO implement me
	panic("implement me")
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't generate hash for given password")
	}
	return string(hash), nil
}

func New() Service {
	return &service{}
}
