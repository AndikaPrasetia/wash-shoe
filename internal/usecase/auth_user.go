// Package usecase: business logics
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/wash-shoe/internal/dto"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/AndikaPrasetia/wash-shoe/internal/redis"
	"github.com/AndikaPrasetia/wash-shoe/internal/repository"
	"github.com/AndikaPrasetia/wash-shoe/internal/sqlc/user"
	utils "github.com/AndikaPrasetia/wash-shoe/internal/utils/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrPasswordMismatch   = errors.New("password and confirm password mismatch")
	ErrEmailAlreadyExists = errors.New("email already exist")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrTokenNotFound      = errors.New("token not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenRevoked       = errors.New("token has been revoked")
)

// AuthUserUsecase defines business logic for auth users
type AuthUserUsecase interface {
	Signup(ctx context.Context, req dto.SignupRequest) (model.AuthUser, string, string, error)
	GetByEmail(ctx context.Context, email string) (*model.AuthUser, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, string, error)
	Logout(ctx context.Context, userID string) error
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}

type authUserUsecase struct {
	authRepo repository.AuthUserRepo
	userRepo repository.UserRepo
	redisCli *redis.RedisClient
}

// NewAuthUserUsecase creates a new AuthUserUsecase
func NewAuthUserUsecase(authRepo repository.AuthUserRepo, userRepo repository.UserRepo, redisCli *redis.RedisClient) AuthUserUsecase {
	return &authUserUsecase{authRepo: authRepo, userRepo: userRepo, redisCli: redisCli}
}

// Signup handles user signup: validates input, creates auth+public user,
// issues tokens, stores refresh token, and logs the action.
func (uc *authUserUsecase) Signup(ctx context.Context, req dto.SignupRequest) (model.AuthUser, string, string, error) {
	// 1. Validate passwords
	if req.Password != req.ConfirmPassword {
		return model.AuthUser{}, "", "", ErrPasswordMismatch
	}
	// 2. Check if email already exists
	existing, err := uc.authRepo.GetAuthUserByEmail(ctx, req.Email)
	if err != nil {
		// Gunakan error dari repository langsung
		if !errors.Is(err, repository.ErrUserNotFound) {
			return model.AuthUser{}, "", "", fmt.Errorf("error checking existing user: %w", err)
		}
		// User tidak ditemukan = boleh lanjut Signup
	} else if existing != nil {
		return model.AuthUser{}, "", "", ErrEmailAlreadyExists
	}
	// jika existing != nil berarti user sudah ada
	if existing != nil {
		return model.AuthUser{}, "", "", ErrEmailAlreadyExists
	}
	// 3. Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return model.AuthUser{}, "", "", fmt.Errorf("hash password: %w", err)
	}
	// 4. Create auth user record (PasswordHash as pgtype.Text)
	authParams := user.CreateAuthUserParams{
		Email:        req.Email,
		PasswordHash: pgtype.Text{String: hash, Valid: true},
	}
	authUser, err := uc.authRepo.Signup(ctx, authParams)
	if err != nil {
		return model.AuthUser{}, "", "", fmt.Errorf("signup auth user: %w", err)
	}
	// 5. Create public user profile
	_, err = uc.userRepo.Create(ctx, user.CreatePublicUserParams{
		ID:          pgtype.UUID{Bytes: uuidFromString(authUser.ID), Valid: true},
		FullName:    req.Username,
		PhoneNumber: pgtype.Text{String: "", Valid: false},
		Role:        "user",
	})
	if err != nil {
		return model.AuthUser{}, "", "", fmt.Errorf("create public user: %w", err)
	}
	// 6. Generate tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(authUser.ID)
	if err != nil {
		return model.AuthUser{}, "", "", fmt.Errorf("generate tokens: %w", err)
	}
	// 7. Store refresh token hash in Redis with expiration
	rtHash := utils.HashToken(refreshToken)
	// Set TTL to refresh token lifetime
	err = uc.redisCli.GetClient().Set(ctx, rtHash, "valid", 7*24*time.Hour).Err() // 7 days
	if err != nil {
		return model.AuthUser{}, "", "", fmt.Errorf("store refresh token in Redis: %w", err)
	}
	// 8. Audit log
	_, _ = uc.authRepo.CreateAuditLog(ctx, user.CreateAuditLogParams{
		ActorID: pgtype.UUID{Bytes: uuidFromString(authUser.ID), Valid: true},
		Action:  "user_register",
		Details: fmt.Appendf(nil, "user %s registered", authUser.Email),
	})
	// 9. Return result
	return authUser, accessToken, refreshToken, nil
}

func (uc *authUserUsecase) GetByEmail(ctx context.Context, email string) (*model.AuthUser, error) {
	return uc.authRepo.GetAuthUserByEmail(ctx, email)
}

// uuidFromString parses a UUID string into uuid.UUID (16-byte array)
func uuidFromString(s string) [16]byte {
	u := uuid.MustParse(s)
	return u
}

func (uc *authUserUsecase) Login(ctx context.Context, req dto.LoginRequest) (string, string, error) {
	// get auth user
	authUser, err := uc.authRepo.GetAuthUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", err
	}

	// validate password
	if !utils.CheckPasswordHash(req.Password, authUser.PasswordHash) {
		// Log failed attempt
		uc.authRepo.CreateAuditLog(ctx, user.CreateAuditLogParams{
			ActorID: pgtype.UUID{Bytes: uuidFromString(authUser.ID), Valid: true},
			Action:  "login_failed",
			Details: fmt.Appendf(nil, "failed login attempt for %s ", req.Email),
		})

		// Return error setelah log
		return "", "", ErrInvalidCredentials
	}

	// get public user data for role
	publicUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", "", err
	}

	// generate tokens
	accessToken, refreshToken, err := utils.GenerateTokenPair(publicUser.ID)
	if err != nil {
		return "", "", err
	}

	// Store refresh token hash in Redis with expiration
	rtHash := utils.HashToken(refreshToken)
	// Set TTL to refresh token lifetime
	err = uc.redisCli.GetClient().Set(ctx, rtHash, "valid", 7*24*time.Hour).Err() // 7 days
	if err != nil {
		return "", "", fmt.Errorf("store refresh token in Redis: %w", err)
	}

	// update last login
	err = uc.authRepo.UpdateLastLogin(ctx, pgtype.UUID{Bytes: uuidFromString(authUser.ID), Valid: true})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uc *authUserUsecase) Logout(ctx context.Context, userID string) error {
	// 1. Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	pgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}

	// 2. Revoke semua refresh token user
	err = uc.authRepo.RevokeAllTokens(ctx, pgUUID)
	if err != nil {
		return fmt.Errorf("failed to revoke tokens: %w", err)
	}

	// 3. [OPTIONAL] Tambahkan access token ke cache revoked
	// (Jika Anda mengimplementasikan cache token yang di-revoke)
	// uc.tokenCache.Add(accessToken, uc.cfg.AccessTokenLifeTime)

	// 4. Audit log (PERBAIKAN DI SINI)
	detailsJSON := []byte(`"user logged out and tokens revoked"`) // String JSON valid

	_, err = uc.authRepo.CreateAuditLog(ctx, user.CreateAuditLogParams{
		ActorID: pgUUID,
		Action:  "logout",
		Details: detailsJSON,
	})

	return err
}

// RefreshToken generates new access and refresh tokens using a valid refresh token
func (uc *authUserUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Parse refresh token
	claims, err := utils.ParseToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// Check if token is a refresh token
	if claims.Type != "refresh" {
		return "", "", ErrInvalidToken
	}

	// Check if token is expired
	if claims.ExpiresAt.Before(time.Now()) {
		return "", "", ErrInvalidToken
	}

	// Check if refresh token is blacklisted in Redis
	rtHash := utils.HashToken(refreshToken)
	val, err := uc.redisCli.GetClient().Get(ctx, rtHash).Result()
	if err != nil || val != "valid" {
		return "", "", ErrTokenRevoked
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := utils.GenerateTokenPair(claims.Subject)
	if err != nil {
		return "", "", fmt.Errorf("generate tokens: %w", err)
	}

	// Blacklist the old refresh token
	err = uc.redisCli.GetClient().Del(ctx, rtHash).Err()
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to blacklist old refresh token: %v\n", err)
	}

	// Store new refresh token hash in Redis with expiration
	newRtHash := utils.HashToken(newRefreshToken)
	// Set TTL to refresh token lifetime
	err = uc.redisCli.GetClient().Set(ctx, newRtHash, "valid", 7*24*time.Hour).Err() // 7 days
	if err != nil {
		return "", "", fmt.Errorf("store new refresh token in Redis: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

func (uc *authUserUsecase) Delete(ctx context.Context, id pgtype.UUID) error {
    return uc.authRepo.Delete(ctx, id)
}
