package http

import (
	"backend/internal/config"
	"backend/internal/delivery/http/handler"
	"backend/internal/delivery/http/middleware"
	"backend/internal/usecase"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	auth *usecase.AuthUsecase,
	product *usecase.ProductUsecase,
	cfg *config.Config,
) *gin.Engine {

	r := gin.Default()
	// for development
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authHandler := handler.NewAuthHandler(auth)
	profileHandler := handler.NewProfileHandler()
	productHandler := handler.NewProductHandler(product)

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

		productGroup := api.Group("/products")
		productGroup.Use(middleware.JWT(cfg))
		{
			productGroup.POST("", productHandler.Create)
			productGroup.GET("", productHandler.FindAll)
			productGroup.GET("/:id", productHandler.FindByID)
			productGroup.PUT("/:id", productHandler.Update)
			productGroup.DELETE("/:id", productHandler.Delete)
		}
	}

	return r
}
