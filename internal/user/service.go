package user

import (
	"context"
	"fmt"
	"github.com/POMBNK/restAPI/pkg/logger"
	"golang.org/x/crypto/bcrypt"
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

	pswrd, err := hashPassword(dto.Password)
	if err != nil {
		return "", fmt.Errorf("failled due hashing password error:%w", err)
	}
	dto.Password = pswrd
	uuid, err := s.storage.Create(ctx, user)
	//TODO: Check if errors are correct
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (s service) GetById(ctx context.Context, id string) (User, error) {
	user, err := s.storage.GetById(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s service) GetAll(ctx context.Context) ([]User, error) {
	users, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (s service) Update(ctx context.Context, dto ToUpdateUserDTO) error {
	user := UpdateUserDto(dto)

	pswrd, err := hashPassword(dto.Password)
	if err != nil {
		return fmt.Errorf("failled due hashing password error:%w", err)
	}
	user.PasswordHash = pswrd

	err = s.storage.Update(ctx, user)
	if err != nil {
		return err
	}

	return err
}

func (s service) Delete(ctx context.Context, id string) error {

	err := s.storage.Delete(ctx, id)
	if err != nil {
		return err
	}
	return err
}

func NewService(storage Storage, logs *logger.Logger) Service {
	return &service{
		storage: storage,
		logs:    logs,
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't generate hash for given password")
	}
	return string(hash), nil
}
