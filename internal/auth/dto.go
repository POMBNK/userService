package auth

type ToSignUpUserDTO struct {
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

type ToSignInUserDTO struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

func CreateSignUpUserDto(dto ToSignUpUserDTO) User {
	return User{
		UserName:     dto.UserName,
		PasswordHash: dto.Password,
		Email:        dto.Email,
	}
}

func CreateSignInUserDto(dto ToSignInUserDTO) User {
	return User{
		PasswordHash: dto.Password,
		Email:        dto.Email,
	}
}
