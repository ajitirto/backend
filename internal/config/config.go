package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	JWTSecret     string
	JWTExpireHour time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env: ", err)
	}

	expireHour, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOUR"))
	if err != nil {
		expireHour = 24
	}

	return &Config{
		Port:          os.Getenv("APP_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPass:        os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpireHour: time.Duration(expireHour) * time.Hour,
	}
}
