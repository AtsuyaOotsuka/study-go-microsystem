package test_db_seeder

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UsersSeeder struct {
	Name         string
	Email        string
	Password     string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func GetUsersSeeders(count int, isDelete bool) []UsersSeeder {
	users := make([]UsersSeeder, count)
	for i := 0; i < count; i++ {
		users[i] = UsersSeeder{
			Name:         fmt.Sprintf("Test User %d", i+1),
			Email:        fmt.Sprintf("user%d@example.com", i+1),
			Password:     fmt.Sprintf("password%d", i+1),
			RefreshToken: fmt.Sprintf("refresh_token_%d", i+1),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			DeletedAt: func() gorm.DeletedAt {
				if isDelete {
					return gorm.DeletedAt{Time: time.Now(), Valid: true}
				}
				return gorm.DeletedAt{}
			}(),
		}
	}
	return users
}
