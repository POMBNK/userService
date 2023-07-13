package author

type ToCreateAuthorDTO struct {
	Name    string `json:"name" `
	SurName string `json:"surname"`
}

type ToUpdateAuthorDTO struct {
	ID      string `json:"ID"`
	Name    string `json:"name" `
	SurName string `json:"surname"`
}

func CreateAuthorDto(dto ToCreateAuthorDTO) Author {
	return Author{
		Name:    dto.Name,
		SurName: dto.SurName,
	}
}

func UpdateAuthorDto(dto ToUpdateAuthorDTO) Author {
	return Author{
		Id:      dto.ID,
		Name:    dto.Name,
		SurName: dto.SurName,
	}
}
