package auth

import "context"

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	GetById(ctx context.Context, id string) (User, error)
}
