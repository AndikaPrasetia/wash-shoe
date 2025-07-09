// Package usecase: business logics
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AndikaPrasetia/wash-shoe/internal/dto"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
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
)

// AuthUserUsecase defines business logic for auth users
type AuthUserUsecase interface {
	Register(ctx context.Context, req dto.SignupRequest) (model.AuthUser, string, string, error)
	GetByEmail(ctx context.Context, email string) (*model.AuthUser, error)
	Login(ctx context.Context, req dto.LoginRequest) (string, string, error)
}

type authUserUsecase struct {
	authRepo repository.AuthUserRepo
	userRepo repository.UserRepo
}

// NewAuthUserUsecase creates a new AuthUserUsecase
func NewAuthUserUsecase(authRepo repository.AuthUserRepo, userRepo repository.UserRepo) AuthUserUsecase {
	return &authUserUsecase{authRepo: authRepo, userRepo: userRepo}
}

// Register handles user signup: validates input, creates auth+public user,
// issues tokens, stores refresh token, and logs the action.
func (uc *authUserUsecase) Register(ctx context.Context, req dto.SignupRequest) (model.AuthUser, string, string, error) {
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
		// User tidak ditemukan = boleh lanjut register
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
	authUser, err := uc.authRepo.Register(ctx, authParams)
	if err != nil {
		return model.AuthUser{}, "", "", fmt.Errorf("register auth user: %w", err)
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
	// 7. Store refresh token
	rtHash := utils.HashToken(refreshToken)
	_, err = uc.authRepo.CreateRefreshToken(ctx, user.CreateRefreshTokenParams{
		UserID:    pgtype.UUID{Bytes: uuidFromString(authUser.ID), Valid: true},
		TokenHash: rtHash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		return model.AuthUser{}, "", "", fmt.Errorf("store refresh token: %w", err)
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

	// save refresh token
	rtHash := utils.HashToken(refreshToken)
	_, err = uc.authRepo.CreateRefreshToken(ctx, user.CreateRefreshTokenParams{
		UserID:    pgtype.UUID{Bytes: uuidFromString(authUser.ID), Valid: true},
		TokenHash: rtHash,
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(30 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		return "", "", err
	}

	// update last login
	err = uc.authRepo.UpdateLastLogin(ctx, pgtype.UUID{Bytes: uuidFromString(authUser.ID), Valid: true})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
