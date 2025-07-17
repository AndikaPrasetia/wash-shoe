// Package handler: routing
package handler

import (
	"errors"
	"net/http"

	"github.com/AndikaPrasetia/wash-shoe/internal/dto"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/AndikaPrasetia/wash-shoe/internal/usecase"
	utils "github.com/AndikaPrasetia/wash-shoe/internal/utils/services"
	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authUC usecase.AuthUserUsecase
	rg     *gin.RouterGroup
}

func NewAuthHandler(authUC usecase.AuthUserUsecase) *authHandler {
	return &authHandler{
		authUC: authUC,
	}
}

func (h *authHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	if err := utils.ValidatePassword(req.Password, 8, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "password_missmatch",
			"message": "Password and confirmation password does not match",
		})
		return
	}

	user, accessToken, refreshToken, err := h.authUC.Signup(c, req)
	if err != nil {
		status := http.StatusInternalServerError
		errType := "server_error"

		switch {
		case errors.Is(err, usecase.ErrEmailAlreadyExists):
			status = http.StatusConflict
			errType = "email_conflict"
		case errors.Is(err, usecase.ErrInvalidCredentials):
			status = http.StatusBadRequest
			errType = "invalid_credential"
		}

		c.JSON(status, gin.H{
			"error":   errType,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})

}
func (h *authHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authUC.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) ||
			errors.Is(err, usecase.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}

func (h *authHandler) Logout(c *gin.Context) {
	// Dapatkan user dari context (diset oleh middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	authUser, ok := user.(model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	// Panggil usecase hanya dengan userID
	err := h.authUC.Logout(c.Request.Context(), authUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.LogoutResponse{
		Message: "Successfully logged out",
	}

	c.JSON(http.StatusOK, response)
}
