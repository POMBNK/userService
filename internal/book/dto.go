package book

import (
	"github.com/POMBNK/restAPI/internal/author"
	"github.com/POMBNK/restAPI/internal/user"
)

type ToCreateBookDTO struct {
	Name     string        `json:"name"`
	UserID   user.User     `json:"userID,omitempty"`
	AuthorID author.Author `json:"authorID"`
}

type ToUpdateBookDTO struct {
	Name     string        `json:"name"`
	UserID   user.User     `json:"userID,omitempty"`
	AuthorID author.Author `json:"authorID"`
}

type ToFindByNameDTO struct {
	Name string `json:"name"`
}

func CreateBookDto(dto ToCreateBookDTO) Book {
	return Book{
		Name:     dto.Name,
		UserID:   dto.UserID,
		AuthorID: dto.AuthorID,
	}
}

func UpdateBookDTO(dto ToUpdateBookDTO) Book {
	return Book{
		Name:     dto.Name,
		UserID:   dto.UserID,
		AuthorID: dto.AuthorID,
	}
}
