// Package controller: routing
package controller

import (
	"errors"
	"net/http"

	"github.com/AndikaPrasetia/wash-shoe/internal/dto"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/AndikaPrasetia/wash-shoe/internal/usecase"
	"github.com/gin-gonic/gin"
)

type authController struct {
	authUC usecase.AuthUserUsecase
	rg     *gin.RouterGroup
}

func NewAuthController(authUC usecase.AuthUserUsecase) *authController {
	return &authController{
		authUC: authUC,
	}
}

func (a *authController) Signup(c *gin.Context) {
	var req dto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "password_missmatch",
			"message": "Password and confirmation password does not match",
		})
	}

	user, accessToken, refreshToken, err := a.authUC.Signup(c, req)
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
func (a *authController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := a.authUC.Login(c.Request.Context(), req)
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

func (a *authController) Logout(c *gin.Context) {
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
	err := a.authUC.Logout(c.Request.Context(), authUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.LogoutResponse{
		Message: "Successfully logged out",
	}

	c.JSON(http.StatusOK, response)
}
