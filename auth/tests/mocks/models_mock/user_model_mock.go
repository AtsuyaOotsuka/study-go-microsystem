package models_mock

import (
	"microservices/auth/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUserMock() *models.User {
	dummyPassword := "password123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dummyPassword), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}

	return &models.User{
		ID:           1,
		Name:         "Test User",
		Email:        "test@example.com",
		Password:     string(passwordHash),
		RefreshToken: "some-refresh-token",
		CreatedAt:    1234567890,
		UpdatedAt:    1234567890,
		DeletedAt:    0,
	}
}
