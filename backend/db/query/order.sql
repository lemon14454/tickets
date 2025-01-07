-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    event_id,
    total_price
) VALUES (
    $1, $2, $3
)
RETURNING *;
