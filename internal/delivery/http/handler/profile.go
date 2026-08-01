package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetUint64("user_id"),
		"email":   c.GetString("email"),
	})
}
