package user

type User struct {
	ID           string `json:"ID" bson:"_id,omitempty"`
	UserName     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
	Email        string `json:"email" bson:"email"`
}

type ToCreateUserDTO struct {
	UserName string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

type ToUpdateUserDTO struct {
	ID       string `json:"ID" bson:"_id,omitempty"`
	UserName string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

func CreateUserDto(dto ToCreateUserDTO) User {
	return User{
		UserName:     dto.UserName,
		PasswordHash: dto.Password,
		Email:        dto.Email,
	}
}

func UpdateUserDto(dto ToUpdateUserDTO) User {
	return User{
		ID:           dto.ID,
		UserName:     dto.UserName,
		PasswordHash: dto.Password,
		Email:        dto.Email,
	}
}
