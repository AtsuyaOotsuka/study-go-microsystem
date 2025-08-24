package jwt

import (
	"errors"
	"microservices/auth/internal/clock_svc"
	"microservices/auth/models"
)

type JwtServiceMockStruct struct {
	Clock clock_svc.ClockInterface
}

func (s *JwtServiceMockStruct) CreateJwt(user *models.User) (string, error) {
	return "mock.jwt.token", nil
}

type JwtServiceFailedMockStruct struct {
	Clock clock_svc.ClockInterface
}

func (s *JwtServiceFailedMockStruct) CreateJwt(user *models.User) (string, error) {
	return "", errors.New("failed to create JWT")
}
