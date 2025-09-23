package jwt_svc

import (
	"fmt"
	"microservices/auth/internal/models"
	"microservices/auth/internal/svc/clock_svc"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtServiceInterface interface {
	CreateJwt(user *models.User) (string, error)
	CreateRefreshToken(c clock_svc.ClockInterface) string
}

type JwtServiceStruct struct {
	Clock  clock_svc.ClockInterface
	Method jwt.SigningMethod
	Key    interface{}
}

func NewJwtService() *JwtServiceStruct {
	return &JwtServiceStruct{
		Clock:  clock_svc.RealClockStruct{},
		Method: jwt.SigningMethodHS256,
		Key:    []byte(os.Getenv("JWT_SECRET")),
	}
}

func (s *JwtServiceStruct) CreateJwt(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   s.Clock.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.Key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *JwtServiceStruct) CreateRefreshToken(c clock_svc.ClockInterface) string {
	refreshToken := fmt.Sprintf("%x", c.Now().UnixNano()) // ランダムな文字列を生成
	return refreshToken
}

func (s *JwtServiceStruct) ValidateJwt(tokenString string) (*models.JwtClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	return &models.JwtClaims{
		UserID: int(claims["sub"].(float64)),
		Email:  claims["email"].(string),
	}, nil
}
