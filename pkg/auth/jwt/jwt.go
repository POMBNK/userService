package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(id string) (string, error) {
	claims := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
