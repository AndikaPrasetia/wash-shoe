package modelutils

import "github.com/golang-jwt/jwt/v5"

type JwtPayloadClaim struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Type   string `json:"type"`
}
