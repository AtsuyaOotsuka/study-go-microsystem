package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUser(
	name string,
	email string,
	password string,
	refreshToken string,
) *User {
	return &User{
		Name:         name,
		Email:        email,
		Password:     password,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

type User struct {
	ID           uint           `gorm:"primaryKey"`
	Name         string         `gorm:"size:255;index"`
	Email        string         `gorm:"unique"`
	Password     string         `gorm:"size:255"`
	RefreshToken string         `gorm:"size:255"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
