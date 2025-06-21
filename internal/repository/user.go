package repository

import (
	"context"
	"errors"

	"github.com/AndikaPrasetia/wash-shoe/internal/db/user"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	sqlErrNoRows    = pgx.ErrNoRows
	ErrUserNotFound = errors.New("user not found")
)

type UserRepo interface {
	Create(ctx context.Context, arg user.CreatePublicUserParams) (model.User, error)
	FindByID(ctx context.Context, id pgtype.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, arg user.UpdatePublicUserParams) (model.User, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type userRepo struct {
	q user.Querier
}

func NewUserRepo(q user.Querier) UserRepo {
	return &userRepo{q: q}
}

func (r *userRepo) Create(ctx context.Context, arg user.CreatePublicUserParams) (model.User, error) {
	pu, err := r.q.CreatePublicUser(ctx, arg)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:          pu.ID.String(),
		FullName:    pu.FullName,
		PhoneNumber: pu.PhoneNumber.String,
		Role:        pu.Role,
		CreatedAt:   pu.CreatedAt.Time,
		UpdatedAt:   pu.UpdatedAt.Time,
	}, nil
}

func (r *userRepo) FindByID(ctx context.Context, id pgtype.UUID) (*model.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber.String,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt.Time,
		UpdatedAt:   u.UpdatedAt.Time,
	}, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	u, err := r.q.GetPublicUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sqlErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &model.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber.String,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt.Time,
		UpdatedAt:   u.UpdatedAt.Time,
	}, nil
}

func (r *userRepo) Update(ctx context.Context, arg user.UpdatePublicUserParams) (model.User, error) {
	u, err := r.q.UpdatePublicUser(ctx, arg)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber.String,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt.Time,
		UpdatedAt:   u.UpdatedAt.Time,
	}, nil
}

func (r *userRepo) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteAuthUser(ctx, id)
}
