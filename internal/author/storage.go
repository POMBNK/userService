package author

import "context"

type Storage interface {
	Create(ctx context.Context, author Author) (string, error)
	GetById(ctx context.Context, id string) (Author, error)
	GetAll(ctx context.Context) ([]Author, error)
	Update(ctx context.Context, author Author) error
	Delete(ctx context.Context, id string) error
}
