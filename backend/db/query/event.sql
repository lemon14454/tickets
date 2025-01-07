-- name: CreateEvent :one
INSERT INTO events (
    host_id,
    start_at,
    name
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetHostEvent :many
SELECT id, name, status, start_at, updated_at, created_at FROM events
WHERE host_id = $1;

-- name: GetAllEvent :many
SELECT id, name, status, start_at FROM events;

-- name: GetEventByID :one
SELECT name, status, start_at FROM events
WHERE id = $1 limit 1;

-- name: GetEventZone :many
SELECT zone, rows, seats, price FROM event_zones
WHERE event_id = $1;

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

-- name: GetEventZones :many
SELECT * FROM event_zones
WHERE event_id = $1;
