package repository

import (
	"context"

	"github.com/AndikaPrasetia/wash-shoe/internal/db"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepo interface {
	Create(ctx context.Context, u model.User) (model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id pgtype.UUID) (*model.User, error)
	Update(ctx context.Context, u db.User) (model.User, error)
	Delete(ctx context.Context, id pgtype.UUID) error
}

type userRepo struct {
	q *db.Queries
}

func NewUserRepo(q *db.Queries) UserRepo {
	return &userRepo{q: q}
}

func (r *userRepo) Create(ctx context.Context, u model.User) (model.User, error) {
	params := db.CreateUserParams{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
	created, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:        created.ID.String(),
		Name:      created.Name,
		Email:     created.Email,
		Role:      created.Role,
		CreatedAt: created.CreatedAt.Time,
		UpdatedAt: created.UpdatedAt.Time,
	}, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (r *userRepo) FindByID(ctx context.Context, id pgtype.UUID) (*model.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.User{
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (r *userRepo) Update(ctx context.Context, u db.User) (model.User, error) {
	params := db.UpdateUserParams{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
	}
	updated, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		ID:        updated.ID.String(),
		Name:      updated.Name,
		Email:     updated.Email,
		Role:      updated.Role,
		CreatedAt: updated.CreatedAt.Time,
		UpdatedAt: updated.UpdatedAt.Time,
	}, nil
}

func (r *userRepo) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteUser(ctx, id)
}
