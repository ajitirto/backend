package http

import (
	"backend/internal/config"
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	auth *usecase.AuthUsecase,
	cfg *config.Config,
) *gin.Engine {

	r := gin.Default()

	authHandler := handler.NewAuthHandler(auth)
	profileHandler := handler.NewProfileHandler()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/signup", authHandler.Signup)
		}

		profileGroup := api.Group("/profile")
		profileGroup.Use(middleware.JWT(cfg))
		{
			profileGroup.GET("", profileHandler.Me)
		}
	}

	return r
}
