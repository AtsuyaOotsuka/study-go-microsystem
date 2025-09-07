package models_mock

import (
	"microservices/auth/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		CreatedAt:    time.Unix(1234567890, 0),
		UpdatedAt:    time.Unix(1234567890, 0),
		DeletedAt:    gorm.DeletedAt{Time: time.Unix(0, 0), Valid: false},
	}
}
