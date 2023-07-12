package book

import (
	"github.com/POMBNK/restAPI/internal/author"
	"github.com/POMBNK/restAPI/internal/user"
)

type Book struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	UserID   user.User     `json:"userID,omitempty"`
	AuthorID author.Author `json:"authorID"`
}
