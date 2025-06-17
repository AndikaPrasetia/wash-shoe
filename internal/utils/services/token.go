package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateTokenPair creates a signed JWT access token and refresh token for the given user ID.
// It reads secret keys and expiration durations from environment variables:
// ACCESS_SECRET, REFRESH_SECRET, ACCESS_EXP (in minutes), REFRESH_EXP (in hours).
func GenerateTokenPair(userID string) (accessToken string, refreshToken string, err error) {
	// Load secrets
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if accessSecret == "" || refreshSecret == "" {
		err = fmt.Errorf("missing token secrets in environment")
		return
	}

	// Parse durations
	accessExpMin := 15
	refreshExpHour := 24
	if v := os.Getenv("ACCESS_EXP"); v != "" {
		if d, parseErr := time.ParseDuration(v + "m"); parseErr == nil {
			accessExpMin = int(d.Minutes())
		}
	}
	if v := os.Getenv("REFRESH_EXP"); v != "" {
		if d, parseErr := time.ParseDuration(v + "h"); parseErr == nil {
			refreshExpHour = int(d.Hours())
		}
	}

	// Create access token
	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(time.Duration(accessExpMin) * time.Minute).Unix(),
		"type": "access",
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		err = fmt.Errorf("sign access token: %w", err)
		return
	}

	// Create refresh token
	refreshClaims := jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(time.Duration(refreshExpHour) * time.Hour).Unix(),
		"type": "refresh",
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString([]byte(refreshSecret))
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
