package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessTokenDuration  = 15
	RefreshTokenDuration = 24 * 31
	AccessTokenName      = "Access-Token"
	RefreshTokenName     = "Refresh-Token"
)

type Claims struct {
	ExpireAt time.Time // когда испортится токен
	UserID   string    // номер идентификатора
	IssuedAt time.Time // когда был создан
}

func (c *Claims) ToMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"exp":    jwt.NewNumericDate(c.ExpireAt),
		"id":     c.UserID,
		"iss_at": jwt.NewNumericDate(time.Now()),
	}
}

type Pair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func newAccessTokenClaims(userID string) *Claims {
	claims := newClaims(userID)
	claims.ExpireAt = claims.ExpireAt.Add(AccessTokenDuration * time.Minute)
	return claims
}

func newRefreshTokenClaims(userID string) *Claims {
	claims := newClaims(userID)
	claims.ExpireAt = claims.ExpireAt.Add(RefreshTokenDuration * time.Hour)
	return claims
}

func newClaims(userID string) *Claims {
	return &Claims{
		ExpireAt: time.Now(),
		UserID:   userID,
		IssuedAt: time.Now(),
	}
}
