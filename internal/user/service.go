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

func (s *service) Create(ctx context.Context, dto ToCreateUserDTO) (string, error) {
	user := CreateUserDto(dto)

	pswrd, err := hashPassword(dto.Password)
	if err != nil {
		return "", fmt.Errorf("failled due hashing password error:%w", err)
	}
	user.PasswordHash = pswrd
	uuid, err := s.storage.Create(ctx, user)
	//TODO: Check if errors are correct
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (s *service) GetById(ctx context.Context, id string) (User, error) {
	user, err := s.storage.GetById(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *service) GetAll(ctx context.Context) ([]User, error) {
	users, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) Update(ctx context.Context, dto ToUpdateUserDTO) error {
	var updUser User
	var pswrd string
	s.logs.Debug("get user by uuid")
	user, err := s.GetById(ctx, dto.ID)
	if err != nil {
		return err
	}
	if dto.UserName == "" {
		dto.UserName = user.UserName
	}
	if dto.Email == "" {
		dto.Email = user.Email
	}
	//TODO: Compare old and new password. Change DTO model.
	if dto.Password == "" {
		dto.Password = user.PasswordHash
		pswrd = user.PasswordHash
	} else {
		pswrd, err = hashPassword(dto.Password)
		if err != nil {
			return fmt.Errorf("failled due hashing password error:%w", err)
		}
	}
	updUser = UpdateUserDto(dto)
	updUser.PasswordHash = pswrd
	err = s.storage.Update(ctx, updUser)
	if err != nil {
		return err
	}
	return err
}

func (s *service) Delete(ctx context.Context, id string) error {

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
