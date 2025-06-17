package repository

import (
	"context"

	"github.com/AndikaPrasetia/wash-shoe/internal/db/user"
	"github.com/AndikaPrasetia/wash-shoe/internal/model"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthUserRepo interface {
	Register(ctx context.Context, arg user.CreateAuthUserParams) (model.AuthUser, error)
	Login(ctx context.Context, email, password string) (model.AuthUser, error)
	Logout(ctx context.Context, userID pgtype.UUID) error
	Delete(ctx context.Context, id pgtype.UUID) error
	CreateRefreshToken(ctx context.Context, arg user.CreateRefreshTokenParams) (model.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, id pgtype.UUID) error
	RevokeAllTokens(ctx context.Context, userID pgtype.UUID) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	CreateAuditLog(ctx context.Context, arg user.CreateAuditLogParams) (model.AuditLog, error)
	ListAuditLogs(ctx context.Context, actorID pgtype.UUID) ([]model.AuditLog, error)
	UpdateLastLogin(ctx context.Context, userID pgtype.UUID) error
}

type authUserRepo struct {
	q user.Querier
}

func NewAuthUserRepo(q user.Querier) AuthUserRepo {
	return &authUserRepo{q: q}
}

func (r *authUserRepo) Register(ctx context.Context, arg user.CreateAuthUserParams) (model.AuthUser, error) {
	au, err := r.q.CreateAuthUser(ctx, arg)
	if err != nil {
		return model.AuthUser{}, err
	}
	return model.AuthUser{
		ID:           au.ID.String(),
		Email:        au.Email,
		PasswordHash: au.PasswordHash.String,
		CreatedAt:    au.CreatedAt.Time,
	}, nil
}

func (r *authUserRepo) Login(ctx context.Context, email, password string) (model.AuthUser, error) {
	auth, err := r.q.GetAuthUserByEmail(ctx, email)
	if err != nil {
		return model.AuthUser{}, err
	}
	// You may want to verify password here before returning
	return model.AuthUser{
		ID:           auth.ID.String(),
		Email:        auth.Email,
		PasswordHash: auth.PasswordHash.String,
	}, nil
}

func (r *authUserRepo) Logout(ctx context.Context, userID pgtype.UUID) error {
	return r.q.UpdateAuthUserLastLogin(ctx, userID)
}

func (r *authUserRepo) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteAuthUser(ctx, id)
}

func (r *authUserRepo) CreateRefreshToken(ctx context.Context, arg user.CreateRefreshTokenParams) (model.RefreshToken, error) {
	tkn, err := r.q.CreateRefreshToken(ctx, arg)
	if err != nil {
		return model.RefreshToken{}, err
	}
	return model.RefreshToken{
		ID:        tkn.ID.String(),
		UserID:    tkn.UserID.String(),
		TokenHash: tkn.TokenHash,
		ExpiresAt: tkn.ExpiresAt.Time,
		Revoked:   tkn.Revoked.Bool,
	}, nil
}

func (r *authUserRepo) RevokeRefreshToken(ctx context.Context, id pgtype.UUID) error {
	return r.q.RevokeRefreshToken(ctx, id)
}

func (r *authUserRepo) RevokeAllTokens(ctx context.Context, userID pgtype.UUID) error {
	return r.q.RevokeAllTokensForUser(ctx, userID)
}

func (r *authUserRepo) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	rt, err := r.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	return &model.RefreshToken{
		ID:        rt.ID.String(),
		UserID:    rt.UserID.String(),
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt.Time,
		Revoked:   rt.Revoked.Bool,
	}, nil
}

func (r *authUserRepo) CreateAuditLog(ctx context.Context, arg user.CreateAuditLogParams) (model.AuditLog, error) {
	al, err := r.q.CreateAuditLog(ctx, arg)
	if err != nil {
		return model.AuditLog{}, err
	}
	return model.AuditLog{
		ID:        al.ID,
		ActorID:   al.ActorID.String(),
		Action:    al.Action,
		Details:   string( al.Details ),
		CreatedAt: al.CreatedAt.Time,
	}, nil
}

func (r *authUserRepo) ListAuditLogs(ctx context.Context, actorID pgtype.UUID) ([]model.AuditLog, error) {
	logs, err := r.q.ListAuditLogs(ctx, actorID)
	if err != nil {
		return nil, err
	}
	res := make([]model.AuditLog, len(logs))
	for i, al := range logs {
		res[i] = model.AuditLog{
			ID:        al.ID,
			ActorID:   al.ActorID.String(),
			Action:    al.Action,
			Details:   string( al.Details ),
			CreatedAt: al.CreatedAt.Time,
		}
	}
	return res, nil
}

func (r *authUserRepo) UpdateLastLogin(ctx context.Context, userID pgtype.UUID) error {
	return r.q.UpdateAuthUserLastLogin(ctx, userID)
}
