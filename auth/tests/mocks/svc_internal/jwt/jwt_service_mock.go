package jwt

import (
	"errors"
	"microservices/auth/internal/models"
	"microservices/auth/internal/svc/clock_svc"
)

type JwtServiceMockStruct struct {
	Clock clock_svc.ClockInterface
}

func (s *JwtServiceMockStruct) CreateJwt(user *models.User) (string, error) {
	return "mock.jwt.token", nil
}

func (s *JwtServiceMockStruct) CreateRefreshToken(c clock_svc.ClockInterface) string {
	return "mock.refresh.token"
}

type JwtServiceFailedMockStruct struct {
	Clock clock_svc.ClockInterface
}

func (s *JwtServiceFailedMockStruct) CreateJwt(user *models.User) (string, error) {
	return "", errors.New("failed to create JWT")
}

func (s *JwtServiceFailedMockStruct) CreateRefreshToken(c clock_svc.ClockInterface) string {
	return ""
}
