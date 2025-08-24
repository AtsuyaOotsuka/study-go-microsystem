package repositories

import (
	"fmt"
	"microservices/auth/models"

	"gorm.io/gorm"
)

type UserRepositoryStruct struct {
	Db *gorm.DB
}

func (r *UserRepositoryStruct) GetByEmail(email string) (*models.User, error) {

	var user models.User
	if err := r.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
