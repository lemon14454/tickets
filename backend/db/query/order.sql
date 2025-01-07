-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    event_id,
    total_price
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserOrders :many
SELECT 
    orders.id,
    orders.event_id,
    orders.created_at,
    orders.total_price,
    events.name,
    events.start_at
FROM orders
JOIN events ON orders.event_id = events.id
WHERE orders.user_id = $1;
