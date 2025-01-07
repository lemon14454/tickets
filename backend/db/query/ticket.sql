-- name: GetRowTickets :many
SELECT * FROM tickets
WHERE event_id = $1
AND zone_id = $2
AND row = $3;

-- name: UpdateTicketsUser :exec
UPDATE tickets SET user_id = sqlc.arg(user_id), order_id = sqlc.arg(order_id)
WHERE id = ANY(sqlc.arg(id)::bigint[]);

-- name: GetTicketsForUpdate :many
SELECT
    tickets.id,
    tickets.user_id,
    tickets.order_id,
    tickets.event_id,
    event_zones.zone,
    tickets.row,
    tickets.seat,
    event_zones.price
FROM tickets 
JOIN event_zones on tickets.zone_id = event_zones.id
WHERE tickets.id = ANY($1::bigint[])
AND tickets.user_id IS NULL
AND tickets.order_id IS NULL
FOR UPDATE;

-- name: GetOrderDetail :many
SELECT
    event_zones.zone,
    tickets.row,
    tickets.seat,
    event_zones.price
FROM tickets 
JOIN event_zones on tickets.zone_id = event_zones.id
WHERE tickets.order_id = $1;
