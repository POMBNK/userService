package user

type User struct {
	ID           string `json:"ID" bson:"_id,omitempty"`
	UserName     string `json:"username" bson:"username"`
	PasswordHash string `json:"-" bson:"password"`
	Email        string `json:"email" bson:"email"`
}

type UserDTO struct {
	UserName     string `json:"username" `
	PasswordHash string `json:"password"`
	Email        string `json:"email"`
}
