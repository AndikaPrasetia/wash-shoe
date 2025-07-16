// Package modelutils: Jwt model
package modelutils

import "github.com/golang-jwt/jwt/v5"

type JwtPayloadClaim struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Role   string `json:"role"` // user, admin
	Type   string `json:"type"` // access, refresh
	Email  string `json:"email"`
}
