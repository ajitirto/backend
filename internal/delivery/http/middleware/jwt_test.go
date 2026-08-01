package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend/internal/config"
	"backend/internal/delivery/http/middleware"
	"backend/internal/infrastructure"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()

	r.Use(middleware.JWT(cfg))

	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": c.GetUint64("user_id"),
			"email":   c.GetString("email"),
		})
	})

	return r
}

func TestJWT_MissingAuthorizationHeader(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "secret",
	}

	router := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{
		"message":"missing authorization header"
	}`, rec.Body.String())
}

func TestJWT_InvalidAuthorizationHeader(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "secret",
	}

	router := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic abc123")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{
		"message":"invalid authorization header"
	}`, rec.Body.String())
}

func TestJWT_InvalidToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret: "secret",
	}

	router := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{
		"message":"invalid token"
	}`, rec.Body.String())
}

func TestJWT_Success(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "super-secret-key",
		JWTExpireHour: time.Duration(24) * time.Hour,
	}

	token, err := infrastructure.GenerateJWT(
		1,
		"admin@example.com",
		cfg.JWTSecret,
		cfg.JWTExpireHour,
	)
	assert.NoError(t, err)

	router := setupRouter(cfg)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	assert.JSONEq(t, `{
		"user_id":1,
		"email":"admin@example.com"
	}`, rec.Body.String())
}
