-- name: CreateOrder :one
INSERT INTO "orders" (user_id, service_type, status)
VALUES ($1, $2, $3)
RETURNING id, user_id, service_type, status, created_at;

-- name: GetOrderByID :one
SELECT id, user_id, service_type, status, created_at
FROM "orders"
WHERE id = $1;

-- name: ListOrdersByUser :many
SELECT id, user_id, service_type, status, created_at
FROM "orders"
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE "orders"
SET status = $2
WHERE id = $1
RETURNING id, user_id, service_type, status, created_at;

-- name: DeleteOrder :exec
DELETE FROM "orders" WHERE id = $1;
