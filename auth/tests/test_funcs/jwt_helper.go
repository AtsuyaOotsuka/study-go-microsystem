package test_funcs

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JwtInfoStruct struct {
	UserID int
	Email  string
}

func JwtConvert(jwtToken string) (JwtInfoStruct, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		fmt.Println("jwtSecret:", string([]byte(os.Getenv("JWT_SECRET"))))
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return JwtInfoStruct{}, err
	}
	return JwtInfoStruct{
		UserID: int(claims["sub"].(float64)),
		Email:  claims["email"].(string),
	}, nil
}
