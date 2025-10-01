package test_funcs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateMockJwtToken(
	userID int,
	email string,
	exp time.Time,
	key []byte,
) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
