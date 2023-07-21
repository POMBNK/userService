package book

import "context"

type Storage interface {
	Create(ctx context.Context, book Book, name, surname string) (id string, err error)
	GetByID(ctx context.Context, id string) (Book, error)
	GetByAuthor(ctx context.Context, author string)
	GetByName(ctx context.Context, userName string)
}
