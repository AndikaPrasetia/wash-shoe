package controller

import (
	"net/http"

	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/gin-gonic/gin"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(c *gin.Context) {
	// Ambil user dari context (diset middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Type assertion ke model.User
	authUser, ok := user.(model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the home page!",
		"user": gin.H{
			"id": authUser.ID,
		},
	})
}
