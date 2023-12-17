package authentication

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	UserID string `json:"userid"`
	jwt.RegisteredClaims
}

const (
	TokenExp  = time.Hour * 3
	SecretKey = "SuperSecretKey"
)

func BuildJWTString(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
