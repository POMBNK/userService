package book

type ToCreateBookDTO struct {
	Name     string `json:"name"`
	UserID   string `json:"userID,omitempty"`
	AuthorID string `json:"authorID"`
}
