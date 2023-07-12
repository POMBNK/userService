package user

type ToCreateUserDTO struct {
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

type ToUpdateUserDTO struct {
	ID       string `json:"ID,omitempty" bson:"_id,omitempty"`
	UserName string `json:"username,omitempty" bson:"username"`
	Password string `json:"password,omitempty" bson:"password"`
	Email    string `json:"email,omitempty" bson:"email"`
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
