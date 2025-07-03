package usecase

import (
	"context"

	"github.com/AndikaPrasetia/wash-shoe/internal/sqlc/user"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/AndikaPrasetia/wash-shoe/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserUsecase defines business logic for public user profiles
type UserUsecase interface {
	Create(ctx context.Context, params user.CreatePublicUserParams) (model.User, error)
	GetByID(ctx context.Context, id pgtype.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, params user.UpdatePublicUserParams) (model.User, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type userUsecase struct {
	repo repository.UserRepo
}

// NewUserUsecase creates a new UserUsecase
func NewUserUsecase(repo repository.UserRepo) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) Create(ctx context.Context, params user.CreatePublicUserParams) (model.User, error) {
	return uc.repo.Create(ctx, params)
}

func (uc *userUsecase) GetByID(ctx context.Context, id pgtype.UUID) (*model.User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *userUsecase) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return uc.repo.FindByEmail(ctx, email)
}

func (uc *userUsecase) Update(ctx context.Context, params user.UpdatePublicUserParams) (model.User, error) {
	return uc.repo.Update(ctx, params)
}

func (uc *userUsecase) Delete(ctx context.Context, id pgtype.UUID) error {
	return uc.repo.Delete(ctx, id)
}
