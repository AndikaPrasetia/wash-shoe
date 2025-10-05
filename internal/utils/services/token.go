package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/AndikaPrasetia/wash-shoe/internal/utils/model-utils"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateTokenPair creates a signed JWT access token and refresh token for the given user ID.
// It reads secret key and expiration durations from environment variables:
// JWT_SECRET, ACCESS_TOKEN_EXP (minutes), REFRESH_TOKEN_EXP (minutes)
func GenerateTokenPair(userID string) (accessToken string, refreshToken string, err error) {
	// Load secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		err = errors.New("JWT_SECRET not set in environment")
		return
	}

	// Parse durations with default values
	accessExpMin := 15     // default 15 minutes
	refreshExpMin := 10080 // default 7 days (10080 minutes)

	if v := os.Getenv("ACCESS_TOKEN_EXP"); v != "" {
		if exp, parseErr := strconv.Atoi(v); parseErr == nil {
			accessExpMin = exp
		}
	}

	if v := os.Getenv("REFRESH_TOKEN_EXP"); v != "" {
		if exp, parseErr := strconv.Atoi(v); parseErr == nil {
			refreshExpMin = exp
		}
	}

	// Create access token
	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(time.Duration(accessExpMin) * time.Minute).Unix(),
		"type": "access",
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString([]byte(secret))
	if err != nil {
		err = fmt.Errorf("sign access token: %w", err)
		return
	}

	// Create refresh token
	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(time.Duration(refreshExpMin) * time.Minute).Unix(),
		"type": "refresh",
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString([]byte(secret))
	if err != nil {
		err = fmt.Errorf("sign refresh token: %w", err)
		return
	}

	return
}

// HashToken returns the SHA-256 hash of the given token in hex encoding.
func HashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

// ParseToken parses a JWT token and returns its claims.
func ParseToken(tokenString string) (modelutils.JwtPayloadClaim, error) {
	// Load secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return modelutils.JwtPayloadClaim{}, errors.New("JWT_SECRET not set in environment")
	}

	// Parse token with claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		&modelutils.JwtPayloadClaim{},
		func(token *jwt.Token) (any, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)

	// Handle parsing error
	if err != nil {
		return modelutils.JwtPayloadClaim{}, fmt.Errorf("token parsing failed: %w", err)
	}

	// Validate token
	if !token.Valid {
		return modelutils.JwtPayloadClaim{}, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*modelutils.JwtPayloadClaim)
	if !ok {
		return modelutils.JwtPayloadClaim{}, errors.New("invalid token claims")
	}

	return *claims, nil
}
