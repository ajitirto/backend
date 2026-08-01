package infrastructure

import (
	"fmt"
	"log"

	"backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}
