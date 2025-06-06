package usecase

import (
	"context"
	// "github.com/AndikaPrasetia/wash-shoe/internal/db"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/AndikaPrasetia/wash-shoe/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserUsecase interface {
	Register(ctx context.Context, u model.User) (model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id pgtype.UUID) (*model.User, error)
	Modify(ctx context.Context, u model.User) (model.User, error)
	Remove(ctx context.Context, id pgtype.UUID) error
}

type userUsecase struct {
	repo repository.UserRepo
}

func NewUserUsecase(repo repository.UserRepo) UserUsecase {
	return &userUsecase{repo: repo}
}

func (uc *userUsecase) Register(ctx context.Context, u model.User) (model.User, error) {
	return uc.repo.Create(ctx, u)
}

func (uc *userUsecase) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return uc.repo.FindByEmail(ctx, email)
}

func (uc *userUsecase) GetByID(ctx context.Context, id pgtype.UUID) (*model.User, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *userUsecase) Modify(ctx context.Context, u model.User) (model.User, error) {
	return uc.repo.Update(ctx, u)
}

func (uc *userUsecase) Remove(ctx context.Context, id pgtype.UUID) error {
	return uc.repo.Delete(ctx, id)
}
