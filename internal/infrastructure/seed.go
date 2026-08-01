package infrastructure

import (
	"backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	// Seed Users
	var userCount int64

	if err := db.Model(&domain.User{}).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount == 0 {
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

		if err := db.Create(&users).Error; err != nil {
			return err
		}
	}

	// Seed Products
	var productCount int64

	if err := db.Model(&domain.Product{}).Count(&productCount).Error; err != nil {
		return err
	}

	if productCount == 0 {
		products := []domain.Product{
			{
				Name:        "MacBook Pro M4",
				Description: "Apple MacBook Pro 14-inch M4",
				Price:       32999000,
				Stock:       10,
			},
			{
				Name:        "Mechanical Keyboard",
				Description: "75% Wireless Mechanical Keyboard",
				Price:       1299000,
				Stock:       25,
			},
			{
				Name:        "Logitech MX Master 3S",
				Description: "Wireless Productivity Mouse",
				Price:       1499000,
				Stock:       15,
			},
			{
				Name:        "27-inch Monitor",
				Description: "4K IPS Monitor",
				Price:       5499000,
				Stock:       8,
			},
			{
				Name:        "USB-C Hub",
				Description: "8-in-1 USB-C Hub",
				Price:       499000,
				Stock:       50,
			},
		}

		if err := db.Create(&products).Error; err != nil {
			return err
		}
	}

	return nil
}
