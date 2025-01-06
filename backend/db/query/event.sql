-- name: CreateEvent :one
INSERT INTO events (
    host_id,
    name
) VALUES (
    $1, $2
)
RETURNING *;

-- name: CreateEventZone :one
INSERT INTO event_zones (
    zone,
    rows,
    seats,
    event_id,
    price
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;
