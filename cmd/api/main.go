package main

import (
	"log"

	"backend/internal/config"
	"backend/internal/delivery/http"
	"backend/internal/infrastructure"
	"backend/internal/repository/postgres"
	"backend/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := infrastructure.NewPostgres(cfg)

	if err != nil {
		log.Fatal(err)
	}

	if err := infrastructure.Migrate(db); err != nil {
		log.Fatal(err)
	}

	if err := infrastructure.Seed(db); err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, cfg)

	router := http.NewRouter(authUsecase, cfg)

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
