-- name: CreateEvent :one
INSERT INTO events (
    host_id,
    start_at,
    name
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetEventByID :one
select name, status, start_at from events
where id = $1 limit 1;

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
select * from event_zones
where event_id = $1;
