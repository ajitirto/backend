package infrastructure

import (
	"backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	var count int64

	if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	password, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []domain.User{
		{
			Name:     "admin",
			Email:    "admin@example.com",
			Password: string(password),
		},
		{
			Name:     "aji",
			Email:    "aji@gmail.com",
			Password: string(password),
		},
	}

	return db.Create(&users).Error
}
