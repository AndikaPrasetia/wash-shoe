package utils

import (
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
	tokenKey := j.cfg.JwtSignaturKey

	claims := modelUtils.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.AppName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTokenLifeTime)),
		},
		UserID: user.ID,
		Role:   user.Role,
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *jwtService) CreateRefreshToken(user model.User) (string, error) {
	tokenKey := j.cfg.JwtSignaturKey

	claims := modelUtils.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.AppName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.RefreshTokenLifeTime)),
		},
		UserID: user.ID,
		Role:   user.Role,
	}

	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
func (j *jwtService) VerifyToken(tokenString string) (modelUtils.JwtPayloadClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(tokenString, &modelUtils.JwtPayloadClaim{}, func(token *jwt.Token) (any, error) {
		return j.cfg.JwtSignaturKey, nil
	})
	if err != nil {
		return modelUtils.JwtPayloadClaim{}, err
	}

	claim, ok := tokenParse.Claims.(*modelUtils.JwtPayloadClaim)
	if !ok {
		return modelUtils.JwtPayloadClaim{}, fmt.Errorf("error claim")
	}

	return *claim, nil
}
