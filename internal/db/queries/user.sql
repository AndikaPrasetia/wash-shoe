-- Auth Users
-- name: CreateAuthUser :one
INSERT INTO auth.users (email, password_hash)
VALUES ($1, $2)
RETURNING *;

-- name: GetAuthUserByEmail :one
SELECT * FROM auth.users WHERE email = $1 LIMIT 1;

-- name: UpdateAuthUserLastLogin :exec
UPDATE auth.users SET last_sign_in_at = NOW() WHERE id = $1;

-- Public Users
-- name: CreatePublicUser :one
INSERT INTO public.users (id, email, full_name, phone_number, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPublicUserByEmail :one
SELECT * FROM public.users WHERE email = $1;

-- name: UpdatePublicUser :one
UPDATE public.users
SET full_name = $2, phone_number = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAuthUser :exec
DELETE FROM auth.users WHERE id = $1;

-- Refresh Tokens
-- name: CreateRefreshToken :one
INSERT INTO auth.refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: RevokeRefreshToken :exec
UPDATE auth.refresh_tokens SET revoked = true WHERE id = $1;

-- Audit Log
-- name: CreateAuditLog :one
INSERT INTO auth.audit_log (actor_id, action, details)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM auth.audit_log WHERE actor_id = $1 ORDER BY created_at DESC;

-- Get User by ID
-- name: GetUserByID :one
SELECT * FROM public.users WHERE id = $1 LIMIT 1;

-- Update User Role (admin only)
-- name: UpdateUserRole :exec
UPDATE public.users SET role = $1 WHERE id = $2;

-- Get Refresh Token by Hash
-- name: GetRefreshTokenByHash :one
SELECT * FROM auth.refresh_tokens 
WHERE token_hash = $1 AND revoked = false;

-- Revoke All Tokens for User
-- name: RevokeAllTokensForUser :exec
UPDATE auth.refresh_tokens 
SET revoked = true 
WHERE user_id = $1;

-- Delete Expired Tokens
-- name: DeleteExpiredTokens :exec
DELETE FROM auth.refresh_tokens 
WHERE expires_at < NOW() AND revoked = true;


