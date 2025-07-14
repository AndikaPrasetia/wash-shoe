package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/wash-shoe/internal/config"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	modelUtils "github.com/AndikaPrasetia/wash-shoe/internal/utils/model-utils"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateAccessToken(user model.User) (string, error)
	CreateRefreshToken(user model.User) (string, error)
	VerifyToken(tokenString string) (modelUtils.JwtPayloadClaim, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{cfg: cfg}
}

func (j *jwtService) CreateAccessToken(user model.User) (string, error) {
	tokenKey := j.cfg.JwtSecretKey

	claims := modelUtils.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.AppName,
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTokenLifeTime)),
		},
		UserID: user.ID,
		Role:   user.Role,
		Type:   "access",
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtService) CreateRefreshToken(user model.User) (string, error) {
	tokenKey := j.cfg.JwtSecretKey

	claims := modelUtils.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.AppName,
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.RefreshTokenLifeTime)),
		},
		UserID: user.ID,
		Role:   user.Role,
		Type:   "refresh",
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
func (j *jwtService) VerifyToken(tokenString string) (modelUtils.JwtPayloadClaim, error) {
	// 1. Parse token dengan claims
	token, err := jwt.ParseWithClaims(
		tokenString,
		&modelUtils.JwtPayloadClaim{},
		func(token *jwt.Token) (any, error) {
			// 2. Validasi metode signing
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.cfg.JwtSecretKey, nil
		},
	)

	// 3. Tangani error parsing
	if err != nil {
		return modelUtils.JwtPayloadClaim{}, fmt.Errorf("token parsing failed: %w", err)
	}

	// 4. Validasi token
	if !token.Valid {
		return modelUtils.JwtPayloadClaim{}, errors.New("invalid token")
	}

	// 5. Ekstrak claims
	claims, ok := token.Claims.(*modelUtils.JwtPayloadClaim)
	if !ok {
		return modelUtils.JwtPayloadClaim{}, errors.New("invalid token claims")
	}

	return *claims, nil
}
