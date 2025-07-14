// Package middleware: app's middleware
package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	utils "github.com/AndikaPrasetia/wash-shoe/internal/utils/services"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Middleware() gin.HandlerFunc
	RequireRole(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService utils.JwtService
}

func NewAuthMiddleware(jwtService utils.JwtService) AuthMiddleware {
	return &authMiddleware{
		jwtService: jwtService,
	}
}

func (a *authMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// [DEBUG] Log request path
		fmt.Printf("Request path: %s\n", c.Request.URL.Path)

		// Ekstrak header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authorization header required",
			})
			return
		}

		// Ekstrak token dari format "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) < 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid authorization format",
			})
			return
		}
		token := tokenParts[1]

		// Verifikasi token
		tokenClaim, err := a.jwtService.VerifyToken(token)
		if err != nil {
			fmt.Printf("Token verification error: %v\n", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			return
		}

		// [DEBUG] Log claims
		fmt.Printf("Token claims: %+v\n", tokenClaim)

		// Pastikan user ID tidak kosong
		userID := tokenClaim.Subject
		if userID == "" {
			// Fallback ke custom claim jika Subject kosong
			userID = tokenClaim.UserID
		}

		if userID == "" {
			fmt.Println("User ID is empty in token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token: missing user ID",
			})
			return
		}

		// Simpan user di context
		c.Set("user", model.User{
			ID:   userID,
			Role: tokenClaim.Role,
		})
		c.Next()
	}
}

func (a *authMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		// check if user role is in the required roles
		userRole := user.(model.User).Role
		if !slices.Contains(roles, userRole) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden access",
			})
			return
		}
		c.Next()
	}
}
