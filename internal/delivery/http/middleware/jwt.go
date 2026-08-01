package middleware

import (
	"net/http"
	"strings"

	"backend/internal/config"
	"backend/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

func JWT(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "missing authorization header",
			})
			return
		}

		const prefix = "Bearer "

		if !strings.HasPrefix(auth, prefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid authorization header",
			})
			return
		}

		token := strings.TrimPrefix(auth, prefix)

		claims, err := infrastructure.ParseJWT(token, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
