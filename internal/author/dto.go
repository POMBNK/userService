package author

type ToCreateAuthorDTO struct {
	Name    string `json:"name" `
	SurName string `json:"surname"`
}

type ToUpdateAuthorDTO struct {
	Name    string `json:"name" `
	SurName string `json:"surname"`
}
