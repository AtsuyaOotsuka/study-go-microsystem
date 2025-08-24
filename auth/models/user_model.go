package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:255;index"`
	Email        string `gorm:"unique"`
	Password     string `gorm:"size:255"`
	RefreshToken string `gorm:"size:255"`
	CreatedAt    int64  `gorm:"autoCreateTime"`
	UpdatedAt    int64  `gorm:"autoUpdateTime"`
	DeletedAt    int64  `gorm:"autoDeleteTime"`
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
