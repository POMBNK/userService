package auth

import (
	"context"
	"fmt"
	"github.com/POMBNK/restAPI/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUP(ctx context.Context, dto ToSignUpUserDTO) (string, error)
	SignIN(ctx context.Context, dto ToSignInUserDTO) (User, error)
}

type service struct {
	storage Storage
	logs    *logger.Logger
}

func (s *service) SignUP(ctx context.Context, dto ToSignUpUserDTO) (string, error) {
	user := CreateSignUpUserDto(dto)

	pswrd, err := hashPassword(dto.Password)
	if err != nil {
		return "", fmt.Errorf("failled due hashing password error:%w", err)
	}
	user.PasswordHash = pswrd
	uuid, err := s.storage.Create(ctx, user)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (s *service) SignIN(ctx context.Context, dto ToSignInUserDTO) (User, error) {

	user, err := s.storage.GetByEmail(ctx, dto.Email)
	if err != nil {
		return User{}, fmt.Errorf("incorrect email or password")
	}

	if err = validatePassword(dto.Password, user.PasswordHash); err != nil {
		return User{}, err
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't generate hash for given password")
	}
	return string(hash), nil
}

func validatePassword(password, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("Incorrect password")
	}
	return nil
}
func NewService(storage Storage, logs *logger.Logger) Service {
	return &service{
		storage: storage,
		logs:    logs,
	}
}
