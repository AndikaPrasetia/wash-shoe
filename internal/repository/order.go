package repository

import (
	"context"

	"github.com/AndikaPrasetia/wash-shoe/internal/db/order"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	// "github.com/jackc/pgx/v5/pgtype"
)

type OrderRepo interface {
	Create(ctx context.Context, u order.Order) (model.Order, error)
	// FindByEmail(ctx context.Context, email string) (*model.Order, error)
	// FindByID(ctx context.Context, id pgtype.UUID) (*model.Order, error)
	// Update(ctx context.Context, u order.Order) (model.Order, error)
	// Delete(ctx context.Context, id pgtype.UUID) error
}

type orderRepo struct {
	q *order.Queries
}

func NewOrderRepo(q *order.Queries) OrderRepo {
	return &orderRepo{q: q}
}

func (r *orderRepo) Create(ctx context.Context, u order.Order) (model.Order, error) {
	params := order.CreateOrderParams{
		UserID:      u.UserID,
		ServiceType: u.ServiceType,
		Status:      u.Status,
	}
	created, err := r.q.CreateOrder(ctx, params)
	if err != nil {
		return model.Order{}, err
	}
	return model.Order{
		UserID:      created.UserID.String(),
		ServiceType: created.ServiceType.String,
		Status:      created.Status.String,
	}, nil
}

// func (r *orderRepo) FindByEmail(ctx context.Context, id string) (*model.Order, error) {
// 	order, err := r.q.GetOrderByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &model.Order{
// 		ID:        user.ID.String(),
// 		Name:      user.Name,
// 		Email:     user.Email,
// 		Role:      user.Role,
// 		CreatedAt: user.CreatedAt.Time,
// 		UpdatedAt: user.UpdatedAt.Time,
// 	}, nil
// }
//
// func (r *userRepo) FindByID(ctx context.Context, id pgtype.UUID) (*model.User, error) {
// 	user, err := r.q.GetUserByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &model.User{
// 		ID:        user.ID.String(),
// 		Name:      user.Name,
// 		Email:     user.Email,
// 		Role:      user.Role,
// 		CreatedAt: user.CreatedAt.Time,
// 		UpdatedAt: user.UpdatedAt.Time,
// 	}, nil
// }
//
// func (r *userRepo) Update(ctx context.Context, u order.Order) (model.Order, error) {
// 	params := order.UpdateOrderStatusParams{
// 		ID:       u.ID,
// 		Name:     u.Name,
// 		Email:    u.Email,
// 		Password: u.Password,
// 		Role:     u.Role,
// 	}
// 	updated, err := r.q.UpdateUser(ctx, params)
// 	if err != nil {
// 		return model.Order{}, err
// 	}
// 	return model.Order{
// 		ID:        updated.ID.String(),
// 		Name:      updated.Name,
// 		Email:     updated.Email,
// 		Role:      updated.Role,
// 		CreatedAt: updated.CreatedAt.Time,
// 		UpdatedAt: updated.UpdatedAt.Time,
// 	}, nil
// }
//
// func (r *userRepo) Delete(ctx context.Context, id pgtype.UUID) error {
// 	return r.q.DeleteUser(ctx, id)
// }
