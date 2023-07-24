package book

import "context"

type Storage interface {
	Create(ctx context.Context, book Book) (id string, err error)
	GetByID(ctx context.Context, id string) (Book, error)
	GetByAuthor(ctx context.Context, author string) ([]Book, error)
	GetByName(ctx context.Context, name string) ([]Book, error)
}
