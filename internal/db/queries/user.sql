-- name: CreateUser :one
INSERT INTO users (name, email, password, role)
VALUES ($1, $2, $3, $4)
RETURNING id, name, email, role, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, name, email, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, name, email, role, created_at, updated_at
FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET name = $2, email = $3, password = $4, role = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, role, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
