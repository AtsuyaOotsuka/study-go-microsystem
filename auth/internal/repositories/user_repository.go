package repositories

import (
	"errors"
	"fmt"
	"microservices/auth/internal/models"

	"gorm.io/gorm"
)

type UserRepositoryStruct struct {
	Db *gorm.DB
}

func (r *UserRepositoryStruct) GetByEmail(email string) (*models.User, error) {

	var user models.User
	if err := r.Db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepositoryStruct) GetByRefreshToken(refreshToken string) (*models.User, error) {
	var user models.User
	if err := r.Db.Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by refresh token: %w", err)
	}

	return &user, nil
}
