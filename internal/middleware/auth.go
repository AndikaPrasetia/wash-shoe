// Package middleware: app's middleware
package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/AndikaPrasetia/wash-shoe/internal/dto"
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
        var aH dto.AuthHeader

        err := c.ShouldBindHeader(&aH)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "message": "Unauthorized",
            })
            return
        }

        token := strings.Replace(aH.AuthorizationHeader, "Bearer ", "", 1)
        tokenClaim, err := a.jwtService.VerifyToken(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "message": "Unauthorized",
            })
            return
        }

        c.Set("user", model.User{
            ID: tokenClaim.UserID,
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
